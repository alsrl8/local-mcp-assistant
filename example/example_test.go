package example

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGreetTool(t *testing.T) {
	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "test-server", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

	// tools 추가
	ss, _ := server.Connect(ctx, st, nil)
	defer func(ss *mcp.ServerSession) {
		_ = ss.Close()
	}(ss)

	client := mcp.NewClient(&mcp.Implementation{Name: "test"}, nil)
	cs, _ := client.Connect(ctx, ct, nil)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "john"},
	})

	// assert 결과
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("response: %v", res)
}
