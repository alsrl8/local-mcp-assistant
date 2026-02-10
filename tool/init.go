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
}
