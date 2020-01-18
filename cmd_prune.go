package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CmdPrune is `direnv prune`
var CmdPrune = &Cmd{
	Name: "prune",
	Desc: "removes old allowed files",
	Action: actionWithConfig(func(env Env, args []string, config *Config) (err error) {
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
				if envrc, err = ioutil.ReadFile(filename); err != nil {
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
	}),
}
