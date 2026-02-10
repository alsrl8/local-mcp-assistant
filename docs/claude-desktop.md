# Claude Desktop

## Configure

### claude_desktop_config.json

```json
{
  "globalShortcut": "",
  "mcpServers": {
    "assistant": {
      "command": "wsl.exe",
      "args": ["/home/<username>/.local/bin/<app>"],
      "env": {
      }
    }
  },
  "preferences": {
    "menuBarEnabled": false,
    "legacyQuickEntryEnabled": false,
    "coworkScheduledTasksEnabled": false,
    "sidebarMode": "chat"
  }
}
```