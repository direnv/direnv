package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

const (
	backSlash   = '\\'
	newLine     = '\n'
	doubleQuote = '"'
)

func printRune(w *bufio.Writer, r rune) {
	switch r {
	case backSlash:
		_, _ = w.WriteRune(backSlash)
		_, _ = w.WriteRune(backSlash)
	case newLine:
		_, _ = w.WriteString("\\n\" +\n\t\"")
	case doubleQuote:
		_, _ = w.WriteRune(backSlash)
		_, _ = w.WriteRune(doubleQuote)
	default:
		if !isASCII(r) {
			panic("only ASCII is supported")
		}
		_, _ = w.WriteRune(r)
	}
}

func isASCII(r rune) bool {
	return r < unicode.MaxASCII
}

func main() {
	flag.Parse()
	packageName := flag.Arg(0)
	constantName := flag.Arg(1)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fprintf(out, "package %s\n\n// %s ...\nconst %s = \"", packageName, constantName, constantName)

	for {
		r, _, err := in.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		printRune(out, r)
	}
	fmt.Fprint(out, "\"\n")
}
