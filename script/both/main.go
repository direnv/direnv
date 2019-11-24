// a little script that outputs stdin back to both stdout and stderr
package main

import (
	"io"
	"os"
)

func main() {
	buf := make([]byte, 2048)

	for {
		size, err := os.Stdin.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			return
		}

		os.Stdout.Write(buf[0:size])
		os.Stderr.Write(buf[0:size])
	}
}
