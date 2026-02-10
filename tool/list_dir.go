package tool

import (
	"context"
	"fmt"
	"local-mcp-assistant/schema"
	"local-mcp-assistant/utils"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ListDir(ctx context.Context, req *mcp.CallToolRequest, input schema.ListDirInput) (
	*mcp.CallToolResult, schema.ListDirOutput, error,
) {
	path := utils.ToWSLPath(input.Path)
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, schema.ListDirOutput{}, fmt.Errorf("failed to read directory: %w", err)
	}

	var list []schema.FileEntry
	for _, e := range entries {
		info, _ := e.Info()
		var size int64
		if info != nil {
			size = info.Size()
		}
		list = append(list, schema.FileEntry{
			Name:  e.Name(),
			IsDir: e.IsDir(),
			Size:  size,
		})
	}
	return nil, schema.ListDirOutput{Entries: list}, nil
}
