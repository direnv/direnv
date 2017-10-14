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
	BackSlash   = '\\'
	NewLine     = '\n'
	DoubleQuote = '"'
)

func printRune(w *bufio.Writer, r rune) {
	switch r {
	case BackSlash:
		w.WriteRune(BackSlash)
		w.WriteRune(BackSlash)
	case NewLine:
		w.WriteString("\\n\" +\n\t\"")
	case DoubleQuote:
		w.WriteRune(BackSlash)
		w.WriteRune(DoubleQuote)
	default:
		if !IsASCII(r) {
			panic("only ASCII is supported")
		}
		w.WriteRune(r)
	}
}

func IsASCII(r rune) bool {
	return r < unicode.MaxASCII
}

func main() {
	flag.Parse()
	packageName := flag.Arg(0)
	constantName := flag.Arg(1)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fprintf(out, "package %s\n\nconst %s = \"", packageName, constantName)

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
