package schema

type QueryDBInput struct {
	Database string `json:"database" jsonschema:"database name from config"`
	Query    string `json:"query" jsonschema:"sql query to execute"`
}

type QueryDBOutput struct {
	Rows []map[string]any `json:"result"`
}

type DBSchemaInput struct {
	Database string `json:"database" jsonschema:"database name from config"`
	Table    string `json:"table,omitempty" jsonschema:"specific table name. if empty, returns all tables"`
}

type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type TableSchema struct {
	Name    string       `json:"name"`
	Columns []ColumnInfo `json:"columns"`
}

type DBSchemaOutput struct {
	Tables []TableSchema `json:"tables"`
}

type ListDBOutput struct {
	Databases []string `json:"databases"`
}
