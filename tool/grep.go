package tool

import (
	"context"
	"errors"
	"fmt"
	"local-mcp-assistant/schema"
	"local-mcp-assistant/utils"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Grep(ctx context.Context, req *mcp.CallToolRequest, input schema.GrepInput) (
	*mcp.CallToolResult, schema.GrepOutput, error,
) {
	path := utils.ToWSLPath(input.Path)
	ctxLines := input.Context
	if ctxLines == 0 {
		ctxLines = 3
	}

	args := []string{
		"-rn",
		fmt.Sprintf("-C%d", ctxLines),
		"--color=never",
		input.Pattern,
		path,
	}
	out, err := exec.CommandContext(ctx, "grep", args...).CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return nil, schema.GrepOutput{Result: "no matches found"}, nil
		}
		return nil, schema.GrepOutput{}, fmt.Errorf("grep failed: %w", err)
	}

	result := string(out)
	const maxLen = 512 * 1024
	if len(result) > maxLen {
		result = result[:maxLen] + "\n... (truncated)"
	}
	return nil, schema.GrepOutput{Result: result}, nil
}
