package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var CmdPrune = &Cmd{
	Name: "prune",
	Desc: "removes old allowed files",
	Fn: func(env Env, args []string) (err error) {
		var config *Config
		var dir *os.File
		var fi os.FileInfo
		var dir_list []string
		var envrc []byte

		if config, err = LoadConfig(env); err != nil {
			return err
		}

		allowed := config.AllowDir()
		if dir, err = os.Open(allowed); err != nil {
			return err
		}
		defer dir.Close()

		if dir_list, err = dir.Readdirnames(0); err != nil {
			return err
		}

		for _, hash := range dir_list {
			filename := path.Join(allowed, hash)
			if fi, err = os.Stat(filename); err != nil {
				return err
			}

			if !fi.IsDir() {
				if envrc, err = ioutil.ReadFile(filename); err != nil {
					return err
				}
				envrc_str := strings.TrimSpace(string(envrc))

				// skip old files, w/o path inside
				if envrc_str == "" {
					continue
				}
				if !fileExists(envrc_str) {
					_ = os.Remove(filename)
				}

			}

		}
		return nil
	},
}
