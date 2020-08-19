package loader

import "strings"

func GetLoadPath(s string) (loadPath string) {
	loadPath = s
	if !strings.HasSuffix(s, "...") {
		if s[len(s)-1] != '/' {
			loadPath += "/"
		}
		loadPath += "..."
	}

	return
}
