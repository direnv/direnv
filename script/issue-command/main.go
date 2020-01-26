package main

import (
	"fmt"
	"os"
	"strings"
)

func escapeData(data string) string {
	data = strings.ReplaceAll(data, "%", "%25")
	data = strings.ReplaceAll(data, "\r", "%0D")
	data = strings.ReplaceAll(data, "\n", "%0A")
	return data
}

func escapeProperty(prop string) string {
	prop = escapeData(prop)
	prop = strings.ReplaceAll(prop, ":", "%3A")
	prop = strings.ReplaceAll(prop, ",", "%2C")
	return prop
}

// Go implementation of
// https://github.com/actions/toolkit/blob/master/packages/core/src/command.ts
func main() {
	var remain []string
	var props []string

	for _, arg := range os.Args[1:] {
		if arg == "--" {
			continue
		}
		if strings.HasPrefix(arg, "--") {
			kv := strings.SplitN(arg[2:], "=", 2)
			if len(kv) != 2 {
				panic(fmt.Sprintf("expected %s to be of form --key=value", arg[2:]))
			}
			props = append(props, fmt.Sprintf("%s=%s", kv[0], escapeProperty(kv[1])))
		} else {
			remain = append(remain, arg)
		}
	}

	if len(remain) != 2 {
		panic(fmt.Sprintf("expected 2 remaining arguments, go %v", remain))
	}

	cmd := remain[0]
	msg := escapeData(remain[1])
	var propStr string
	if len(props) > 0 {
		propStr = " " + strings.Join(props, ",")
	}

	fmt.Printf("::%s%s::%s", cmd, propStr, msg)
}
