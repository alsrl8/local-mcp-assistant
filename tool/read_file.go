package tool

import (
	"context"
	"fmt"
	"local-mcp-assistant/schema"
	"local-mcp-assistant/utils"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ReadFile(ctx context.Context, req *mcp.CallToolRequest, input schema.ReadFileInput) (
	*mcp.CallToolResult, schema.ReadFileOutput, error,
) {
	path := utils.ToWSLPath(input.Path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, schema.ReadFileOutput{}, fmt.Errorf("failed to read file: %w", err)
	}

	// 너무 큰 파일 방어
	const maxSize = 512 * 1024 // 512KB
	if len(data) > maxSize {
		return nil, schema.ReadFileOutput{
			Content: string(data[:maxSize]) + "\n... (truncated)",
		}, nil
	}
	return nil, schema.ReadFileOutput{Content: string(data)}, nil
}
