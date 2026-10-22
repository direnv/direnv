package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/direnv/direnv/v2/gzenv"
)

func TestDumpArg(t *testing.T) {
	dump := gzenv.Marshal(Env{"FOO": "bar"})

	// An inline dump string is returned unchanged.
	got, err := dumpArg(dump)
	if err != nil {
		t.Fatal("unexpected error for inline dump:", err)
	}
	if got != dump {
		t.Error("inline dump was modified")
	}

	// A path to a file containing a dump is read from disk,
	// like `direnv dump > dump` followed by `direnv show_dump dump`.
	path := filepath.Join(t.TempDir(), "dump")
	if err := os.WriteFile(path, []byte(dump+"\n"), 0600); err != nil {
		t.Fatal("could not write dump file:", err)
	}

	got, err = dumpArg(path)
	if err != nil {
		t.Fatal("unexpected error for dump file:", err)
	}

	var env Env
	if err := gzenv.Unmarshal(got, &env); err != nil {
		t.Fatal("dump read from file did not unmarshal:", err)
	}
	if env["FOO"] != "bar" {
		t.Error("FOO != bar", env["FOO"])
	}

	// A directory is not mistaken for a dump file.
	dir := t.TempDir()
	got, err = dumpArg(dir)
	if err != nil {
		t.Fatal("unexpected error for directory:", err)
	}
	if got != dir {
		t.Error("directory argument was modified")
	}
}
