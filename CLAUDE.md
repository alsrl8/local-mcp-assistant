# CLAUDE.md

## Project Overview

**mcp-assistant** — A Go MCP (Model Context Protocol) server that runs on WSL and provides local filesystem access and database query tools to MCP clients (e.g., Claude Desktop). Automatically converts Windows paths to WSL paths.

## Tech Stack

- Go 1.25
- MCP SDK: `github.com/modelcontextprotocol/go-sdk v1.3.0`
- Config: Viper (`github.com/spf13/viper`)
- Database: PostgreSQL (`github.com/lib/pq`)
- Transport: stdio

## Project Structure

```
├── main.go          # Entry point, MCP server setup, stdio transport
├── config/          # Viper-based config loading (database connections)
├── tool/            # Tool implementations (one file per tool)
│   ├── init.go      # Tool registration via mcp.AddTool()
│   ├── list_dir.go  # Directory listing
│   ├── read_file.go # File reading (512KB limit)
│   ├── grep.go      # Regex search with context lines
│   └── database.go  # DB tools: list_databases, describe_schema, query_db, execute_db
├── schema/          # Input/output structs with json + jsonschema tags
├── utils/           # Windows-to-WSL path conversion
├── example/         # Reference tool implementation + test
├── docs/            # Documentation and architecture diagram
├── Makefile         # build, install, deploy targets
└── go.mod / go.sum
```

## Build & Run

```bash
# Build
make build          # Output: bin/mcp-assistant

# Install to ~/.local/bin
make install

# Build + install
make deploy

# Run (requires stdio MCP client)
MCP_ASSISTANT_CONFIG=/path/to/config.yaml ~/.local/bin/mcp-assistant
```

## Test

```bash
go test ./...
```

The `example/` directory contains a reference test using in-memory MCP transports.

## Configuration

Set `MCP_ASSISTANT_CONFIG` env var to a YAML file path. Config defines database connections and writable table permissions:

```yaml
databases:
  - name: mydb
    host: localhost
    port: 5432
    user: postgres
    password: secret
    dbname: myapp
    writable_tables:  # only these tables allow INSERT/UPDATE/DELETE
      - users
      - logs
```

## Logging

Logs go to `/tmp/mcp-assistant.log` using Go's `slog` package.

## Adding a New Tool

1. Define input/output structs in `schema/newtool.go` (use `json` + `jsonschema` struct tags)
2. Implement handler in `tool/newtool.go` following the signature:
   ```go
   func NewTool(ctx context.Context, req *mcp.CallToolRequest, input schema.Input) (
       *mcp.CallToolResult, schema.Output, error,
   )
   ```
3. Register in `tool/init.go` via `mcp.AddTool(server, &mcp.Tool{...}, NewTool)`

## Code Conventions

- Every tool handler converts Windows paths: `path := utils.ToWSLPath(input.Path)`
- Errors wrapped with `fmt.Errorf("context: %w", err)`
- Schema structs use `json:"field"` and `jsonschema:"description"` tags
- Structured logging: `slog.Info/Error/Warn("msg", "key", "val")`

## Key Constraints

- `query_db`: SELECT only
- `execute_db`: INSERT/UPDATE/DELETE only, restricted to `writable_tables` in config
- File reads truncated at 512KB
- Grep output capped at 512KB
- All filesystem paths auto-converted from Windows to WSL format
