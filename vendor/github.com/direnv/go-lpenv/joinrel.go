package lpenv

import "path/filepath"

// Like filepath.Join, but only joins if the file has a relative path
func JoinRel(path, file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(path, file)
}
