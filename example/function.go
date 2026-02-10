package example

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}
