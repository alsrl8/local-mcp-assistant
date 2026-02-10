package utils

import "strings"

func ToWSLPath(path string) string {
	if len(path) >= 2 && path[1] == ':' {
		drive := strings.ToLower(string(path[0]))
		rest := strings.ReplaceAll(path[2:], "\\", "/")
		return "/mnt/" + drive + rest
	}
	return path
}
