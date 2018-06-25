package sanitize

import (
	"path/filepath"
	"strings"
)

const BadFilenameChars = "\n\t\r=*?:[]\"<>|(){}&'!;$`\\"

//assumes unicode
func SanitizeFilename(filename string) string {
	filename = strings.Map(func(r rune) rune {
		if r < 32 || (r > 126 && r < 160) || strings.ContainsRune(BadFilenameChars, r) {
			return -1
		}
		return r
	}, filename)

	filename = strings.TrimSpace(filename)
	filename = strings.TrimLeft(filename, " -~")
	filename = strings.TrimSpace(filename)
	filename = strings.Replace(filename, " ", "_", -1)
	filename = strings.Replace(filename, "_-_", "", -1)
	var (
		st   strings.Builder
		prev rune
	)
	st.Grow(len(filename))
	for _, r := range filename {
		switch r {
		case '_':
			if r != prev {
				st.WriteRune(r)
			}
		default:
			st.WriteRune(r)
		}
		prev = r
	}
	return st.String()
}

func SanitizeFilepath(path string) string {
	var newpath string
	for _, filename := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		newpath = filepath.Join(newpath, SanitizeFilename(filename))
	}
	return newpath
}
