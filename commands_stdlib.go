//
// Commands that we want to expose in the stdlib.
// Generally they exist because of cross-platform issues.
//

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Makes a path absolute, relative to another "wd" path
func expandPath(path string, wd string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(wd, path)
}

func ExpandPath(args []string) (err error) {
	var path, wd string

	flagset := flag.NewFlagSet("direnv expand-path", flag.ExitOnError)
	flagset.Parse(args[1:])

	path = flagset.Arg(0)
	if path == "" {
		return fmt.Errorf("PATH missing")
	}

	wd = flagset.Arg(1)
	if wd == "" {
		if wd, err = os.Getwd(); err != nil {
			return
		}
	}

	absPath := expandPath(path, wd)
	fmt.Println(absPath)

	return nil
}

func FileHash(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv file-hash", flag.ExitOnError)
	flagset.Parse(args[1:])

	path := flagset.Arg(0)
	if path == "" {
		return fmt.Errorf("PATH missing")
	}

	fd, err := os.Open(path)
	if err != nil {
		return
	}

	hasher := sha256.New()
	io.Copy(hasher, fd)

	num := hasher.Sum(nil)
	// str := base64.URLEncoding.EncodeToString(num)

	// fmt.Printf("%v\n", str)
	fmt.Printf("%x\n", num)

	return
}

func FileMtime(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv file-mtime", flag.ExitOnError)
	flagset.Parse(args[1:])

	path := flagset.Arg(0)
	if path == "" {
		return fmt.Errorf("PATH missing")
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err, path)
		return
	}

	fmt.Println(fileInfo.ModTime().Unix())
	return
}

func FindUp(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv file-mtime", flag.ExitOnError)
	flagset.Parse(args[1:])

	path := flagset.Arg(0)
	fileName := flagset.Arg(1)

	if fileName == "" {
		fileName = path
		path = ""
	}

	if fileName == "" {
		return fmt.Errorf("Missing PATH and FILE_NAME arguments")
	}

	if path == "" {
		if path, err = os.Getwd(); err != nil {
			return
		}
	}

	newPath := findUp(path, fileName)
	if newPath != "" {
		fmt.Println(newPath)
	}

	return
}

func Stdlib(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv stdlib", flag.ExitOnError)
	flagset.Parse(args[1:])

	fmt.Println(STDLIB)
	return
}
