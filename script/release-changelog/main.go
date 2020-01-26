package main

import (
	"bufio"
	"os"
	"strings"
)

func run() error {
	headers := 0

	f, err := os.Open("CHANGELOG.md")
	if err != nil {
		return err
	}

	r := bufio.NewReader(f)

	prev, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.HasPrefix(line, "======") {
			headers++
		}
		if headers > 1 {
			return nil
		}

		os.Stdout.WriteString(prev)
		prev = line
	}
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
