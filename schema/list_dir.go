package schema

type ListDirInput struct {
	Path string `json:"path" jsonschema:"directory path to list"`
}

type ListDirOutput struct {
	Entries []FileEntry `json:"entries"`
}

type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}
