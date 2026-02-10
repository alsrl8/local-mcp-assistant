package main

import (
	"context"
	"local-mcp-assistant/config"
	"local-mcp-assistant/tool"
	"log"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	f, err := os.OpenFile("/tmp/mcp-assistant.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(f, nil)))

	slog.Info("mcp-assistant started")

	configPath := os.Getenv("MCP_ASSISTANT_CONFIG")
	if configPath == "" {
		slog.Warn("MCP_ASSISTANT_CONFIG not set")
		return
	}

	_, err = config.Load(configPath)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		return
	}
}

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-assistant",
			Version: "v1.0.0",
		},
		nil,
	)

	// add tools
	tool.AddTools(server)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
