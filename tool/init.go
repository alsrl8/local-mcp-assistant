package tool

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_dir",
		Description: "List files and directories with sizes on the user's local machine. " +
			"Use when user asks to see folder contents, directory listings, or file structures. " +
			"Accepts Windows paths (C:\\...) and Unix paths.",
	}, ListDir)

	mcp.AddTool(server, &mcp.Tool{
		Name: "read_file",
		Description: "Read the full content of a file. " +
			"Use when user asks to view, read, or cat a file. " +
			"Accepts Windows and Unix paths. Large files are truncated at 512KB.",
	}, ReadFile)

	mcp.AddTool(server, &mcp.Tool{
		Name: "grep",
		Description: "Search for a pattern in a file or directory using regex. " +
			"Returns matching lines with line numbers and surrounding context. " +
			"Use when user asks to search, find, or grep for text in files.",
	}, Grep)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_databases",
		Description: "List available database connections from config.",
	}, ListDatabases)

	mcp.AddTool(server, &mcp.Tool{
		Name: "describe_schema",
		Description: "Describe database schema. Returns table names and column details. " +
			"Use when user asks about table structure, columns, or database schema.",
	}, DescribeSchema)

	mcp.AddTool(server, &mcp.Tool{
		Name: "query_db",
		Description: "Execute a read-only SQL query (SELECT only) against a database. " +
			"Use list_databases to find available databases and describe_schema to understand table structures first.",
	}, QueryDB)

	mcp.AddTool(server, &mcp.Tool{
		Name: "execute_db",
		Description: "Execute a write SQL query (INSERT, UPDATE, DELETE) against a database. " +
			"Only tables listed in writable_tables config are allowed. " +
			"Use describe_schema to understand table structures first.",
	}, ExecuteDB)
}
