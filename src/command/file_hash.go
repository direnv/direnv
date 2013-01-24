package command

import (
	"crypto/sha256"
	// "encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

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
