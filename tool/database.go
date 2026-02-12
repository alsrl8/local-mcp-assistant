package tool

import (
	"context"
	"database/sql"
	"fmt"
	"local-mcp-assistant/config"
	"local-mcp-assistant/schema"
	"log/slog"
	"strings"

	_ "github.com/lib/pq"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func connectDB(name string) (*sql.DB, error) {
	cfg := config.Get()
	var dbCfg *config.DatabaseConfig
	for _, d := range cfg.Databases {
		if d.Name == name {
			dbCfg = &d
			break
		}
	}
	if dbCfg == nil {
		return nil, fmt.Errorf("database not found: %s", name)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName)
	slog.Info("connecting to database", "dsn", dsn)
	return sql.Open("postgres", dsn)
}

func ListDatabases(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (
	*mcp.CallToolResult, schema.ListDBOutput, error,
) {
	cfg := config.Get()
	var names []string
	for _, d := range cfg.Databases {
		names = append(names, d.Name)
	}
	return nil, schema.ListDBOutput{Databases: names}, nil
}

func QueryDB(ctx context.Context, req *mcp.CallToolRequest, input schema.QueryDBInput) (
	*mcp.CallToolResult, schema.QueryDBOutput, error,
) {
	q := strings.TrimSpace(strings.ToUpper(input.Query))
	if !strings.HasPrefix(q, "SELECT") {
		return nil, schema.QueryDBOutput{}, fmt.Errorf("only SELECT queries are allowed")
	}

	db, err := connectDB(input.Database)
	if err != nil {
		return nil, schema.QueryDBOutput{}, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, err := db.QueryContext(ctx, input.Query)
	if err != nil {
		return nil, schema.QueryDBOutput{}, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	cols, _ := rows.Columns()
	var result []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		_ = rows.Scan(ptrs...)
		row := make(map[string]any)
		for i, col := range cols {
			row[col] = vals[i]
		}
		result = append(result, row)
	}

	return nil, schema.QueryDBOutput{Rows: result}, nil
}

func getWritableTables(dbName string) []string {
	cfg := config.Get()
	for _, d := range cfg.Databases {
		if d.Name == dbName {
			return d.WritableTables
		}
	}
	return nil
}

func isWritableTable(dbName, tableName string) bool {
	for _, t := range getWritableTables(dbName) {
		if strings.EqualFold(t, tableName) {
			return true
		}
	}
	return false
}

// extractTargetTable parses the target table name from INSERT, UPDATE, DELETE queries.
func extractTargetTable(query string) (string, error) {
	tokens := strings.Fields(query)
	if len(tokens) < 3 {
		return "", fmt.Errorf("query too short to determine target table")
	}

	switch strings.ToUpper(tokens[0]) {
	case "INSERT":
		// INSERT INTO <table>
		if strings.ToUpper(tokens[1]) != "INTO" {
			return "", fmt.Errorf("invalid INSERT syntax")
		}
		return strings.ToLower(tokens[2]), nil
	case "UPDATE":
		// UPDATE <table>
		return strings.ToLower(tokens[1]), nil
	case "DELETE":
		// DELETE FROM <table>
		if strings.ToUpper(tokens[1]) != "FROM" {
			return "", fmt.Errorf("invalid DELETE syntax")
		}
		return strings.ToLower(tokens[2]), nil
	default:
		return "", fmt.Errorf("unsupported statement: %s", tokens[0])
	}
}

func ExecuteDB(ctx context.Context, req *mcp.CallToolRequest, input schema.ExecuteDBInput) (
	*mcp.CallToolResult, schema.ExecuteDBOutput, error,
) {
	q := strings.TrimSpace(input.Query)
	upper := strings.ToUpper(q)

	// only INSERT, UPDATE, DELETE allowed
	if !strings.HasPrefix(upper, "INSERT") && !strings.HasPrefix(upper, "UPDATE") && !strings.HasPrefix(upper, "DELETE") {
		return nil, schema.ExecuteDBOutput{}, fmt.Errorf("only INSERT, UPDATE, DELETE queries are allowed. use query_db for SELECT")
	}

	tableName, err := extractTargetTable(q)
	if err != nil {
		return nil, schema.ExecuteDBOutput{}, fmt.Errorf("failed to parse table name: %w", err)
	}

	if !isWritableTable(input.Database, tableName) {
		return nil, schema.ExecuteDBOutput{}, fmt.Errorf("table %q is not in writable_tables for database %q", tableName, input.Database)
	}

	db, err := connectDB(input.Database)
	if err != nil {
		return nil, schema.ExecuteDBOutput{}, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	result, err := db.ExecContext(ctx, input.Query)
	if err != nil {
		return nil, schema.ExecuteDBOutput{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	slog.Info("execute_db", "database", input.Database, "table", tableName, "rows_affected", rowsAffected)

	return nil, schema.ExecuteDBOutput{RowsAffected: rowsAffected}, nil
}

func DescribeSchema(ctx context.Context, req *mcp.CallToolRequest, input schema.DBSchemaInput) (
	*mcp.CallToolResult, schema.DBSchemaOutput, error,
) {
	db, err := connectDB(input.Database)
	if err != nil {
		return nil, schema.DBSchemaOutput{}, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	query := `
		SELECT table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'`
	var args []any

	if input.Table != "" {
		query += " AND table_name = $1"
		args = append(args, input.Table)
	}
	query += " ORDER BY table_name, ordinal_position"

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, schema.DBSchemaOutput{}, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	tableMap := make(map[string]*schema.TableSchema)
	var order []string
	for rows.Next() {
		var tbl, col, dtype, nullable string
		_ = rows.Scan(&tbl, &col, &dtype, &nullable)
		if _, ok := tableMap[tbl]; !ok {
			tableMap[tbl] = &schema.TableSchema{Name: tbl}
			order = append(order, tbl)
		}
		tableMap[tbl].Columns = append(tableMap[tbl].Columns, schema.ColumnInfo{
			Name:     col,
			Type:     dtype,
			Nullable: nullable == "YES",
		})
	}

	var tables []schema.TableSchema
	for _, name := range order {
		tables = append(tables, *tableMap[name])
	}
	return nil, schema.DBSchemaOutput{Tables: tables}, nil
}
