package main

import (
	"path/filepath"
	"fmt"
)

func ExpandPath(args []string) error {
	if len(args) > 1 {
		absPath, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		fmt.Println(absPath)
	}
	return nil
}

