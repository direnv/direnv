package cmd

import (
	"os"
	"path"
	"strings"
)

// CmdPrune is `direnv prune`
var CmdPrune = &Cmd{
	Name:   "prune",
	Desc:   "Removes old allowed files",
	Action: actionWithConfig(cmdPruneAction),
}

func cmdPruneAction(_ Env, _ []string, config *Config) (err error) {
	var dir *os.File
	var fi os.FileInfo
	var dirList []string
	var envrc []byte

	allowed := config.AllowDir()
	if dir, err = os.Open(allowed); err != nil {
		return err
	}
	defer dir.Close()

	if dirList, err = dir.Readdirnames(0); err != nil {
		return err
	}

	for _, hash := range dirList {
		filename := path.Join(allowed, hash)
		if fi, err = os.Stat(filename); err != nil {
			return err
		}

		if !fi.IsDir() {
			if envrc, err = os.ReadFile(filename); err != nil {
				return err
			}
			envrcStr := strings.TrimSpace(string(envrc))

			// skip old files, w/o path inside
			if envrcStr == "" {
				continue
			}
			if !fileExists(envrcStr) {
				_ = os.Remove(filename)
			}

		}

	}
	return nil
}
