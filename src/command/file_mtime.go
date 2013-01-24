package command

import (
	"flag"
	"fmt"
	"os"
)

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
