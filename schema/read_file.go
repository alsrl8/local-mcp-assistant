package schema

type ReadFileInput struct {
	Path string `json:"path" jsonschema:"file path to read"`
}

type ReadFileOutput struct {
	Content string `json:"content"`
}
