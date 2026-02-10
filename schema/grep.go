package schema

type GrepInput struct {
	Pattern string `json:"pattern" jsonschema:"search pattern (regex supported)"`
	Path    string `json:"path" jsonschema:"file or directory path to search"`
	Context int    `json:"context,omitempty" jsonschema:"number of context lines around each match, default 3"`
}

type GrepOutput struct {
	Result string `json:"result"`
}
