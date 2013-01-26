package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Command func(args []string) error

func Diff(args []string) (err error) {
	var reverse bool

	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.BoolVar(&reverse, "reverse", false, "Reverses the diff")
	flagset.Parse(args[1:])

	oldEnvStr := flagset.Arg(0)

	if oldEnvStr == "" {
		return fmt.Errorf("Missing OLD_ENV argument")
	}

	oldEnv, err := ParseEnv(oldEnvStr)
	if err != nil {
		return fmt.Errorf("Parse env error: %v", err)
	}

	newEnv := FilteredEnv()

	var diff EnvDiff
	if reverse {
		diff = DiffEnv(oldEnv, newEnv)
	} else {
		diff = DiffEnv(newEnv, oldEnv)
	}

	fmt.Println(EnvToBash(diff))

	return
}

func Dump(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.Parse(args[1:])

	e := FilteredEnv()
	str, err := e.Serialize()
	if err != nil {
		return
	}
	fmt.Println(str)
	return
}

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

func Stdlib(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv stdlib", flag.ExitOnError)
	flagset.Parse(args[1:])

	fmt.Println(STDLIB)
	return
}

func Export(args []string) (err error) {
	// TODO

	return
}

// NOTE: direnv hook $0
// $0 starts with "-" and go tries to parse it as an argument
func Hook(args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	} else {
		// Try to find out the shell on Linux systems
		ppid := os.Getppid()
		data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", ppid))
		if err != nil {
			return fmt.Errorf("Please specify a target shell")
		}

		target = string(data)
	}

	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	target = filepath.Base(target)

	switch target {
	case "bash":
		fmt.Println("PROMPT_COMMAND=\"eval \\`direnv export\\`;$PROMPT_COMMAND")
	case "zsh":
		fmt.Println("direnv_hook() { eval `direnv export` }; [[ -z $precmd_functions ]] && precmd_functions=(); precmd_functions=($precmd_functions direnv_hook)")
	default:
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	return
}
