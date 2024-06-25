package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// loadTempTestEnv will save a copy of the current environment, as specified in
// the keys of e, and then set the environment to the values of e. A function is
// returned that can be used to restore the environment to its original state,
// useful in a defer function.
func loadTempTestEnv(t *testing.T, e Env) (Env, func()) {
	originalEnv := make(map[string]*string)
	for key, value := range e {
		origValue, exists := os.LookupEnv(key)
		if exists {
			originalEnv[key] = &origValue
		} else {
			originalEnv[key] = nil
		}
		t.Logf("setting temporary env: %s=%s\n", key, value)
		_ = os.Setenv(key, value)
	}

	ne := GetEnv()

	keys := make([]string, 0, len(ne))
	for key := range ne {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// for _, key := range keys {
	// 	t.Logf("loaded: %s=%s\n", key, ne[key])
	// }

	return ne, func() {
		for key, origValue := range originalEnv {
			if origValue != nil {
				t.Logf("unsetting temporary env: %s=%s\n", key, *origValue)
				os.Setenv(key, *origValue)
			} else {
				t.Logf("unsetting temporary env: %s\n", key)
				os.Unsetenv(key)
			}
		}
	}
}

func TestCygpathLookPathOs(t *testing.T) {

	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_unix.go;l=52
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_unix_test.go

	if Cygpath == nil {
		t.Skip(errCygpath.Error())
	}

	driveLetter := os.Getenv("SYSTEMDRIVE")
	if len(driveLetter) == 0 {
		t.Error("SYSTEMDRIVE is empty")
	}

	if !strings.Contains(driveLetter, ":") {
		t.Error("SYSTEMDRIVE expected :", driveLetter)
	}

	pathenv := os.Getenv("PATH")

	if !strings.Contains(pathenv, ";") {
		t.Errorf("expected os PATH to use Windows separator: %s", pathenv)
	}

	if !strings.Contains(pathenv, driveLetter) {
		t.Errorf("expected os PATH to use Windows drive letter: %s", pathenv)
	}

	var err error
	var p string

	p, err = lookPath("hostname", pathenv)
	if err != nil {
		t.Errorf("expected hostname in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("hostname.exe", pathenv)
	if err != nil {
		t.Errorf("expected hostname.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/hostname", pathenv)
	if err != nil {
		t.Errorf("expected hostname in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/hostname.exe", pathenv)
	if err != nil {
		t.Errorf("expected hostname.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("env", pathenv)
	if err != nil {
		t.Errorf("expected env in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("env.exe", pathenv)
	if err != nil {
		t.Errorf("expected env.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/env", pathenv)
	if err != nil {
		t.Errorf("expected env in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/env.exe", pathenv)
	if err != nil {
		t.Errorf("expected env.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}
}

func TestCygpathLookPathEnv(t *testing.T) {

	if Cygpath == nil {
		t.Skip(errCygpath.Error())
	}

	driveLetter := os.Getenv("SYSTEMDRIVE")
	if len(driveLetter) == 0 {
		t.Error("SYSTEMDRIVE is empty")
	}

	if !strings.Contains(driveLetter, ":") {
		t.Error("SYSTEMDRIVE expected :", driveLetter)
	}

	env, restoreEnv := loadTempTestEnv(t, Env{})
	defer restoreEnv()

	pathenv := env["PATH"]

	if strings.Contains(pathenv, ";") {
		t.Errorf("expected env PATH not to use Windows separator: %s", pathenv)
	}

	if strings.Contains(pathenv, driveLetter) {
		t.Errorf("expected env PATH not to use Windows drive letter: %s", pathenv)
	}

	var err error
	var p string

	p, err = lookPath("hostname", pathenv)
	if err != nil {
		t.Errorf("expected hostname in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("hostname.exe", pathenv)
	if err != nil {
		t.Errorf("expected hostname.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/hostname", pathenv)
	if err != nil {
		t.Errorf("expected hostname in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/hostname.exe", pathenv)
	if err != nil {
		t.Errorf("expected hostname.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "hostname.exe") {
		t.Errorf("expected hostname.exe in name: %s", p)
	}

	p, err = lookPath("env", pathenv)
	if err != nil {
		t.Errorf("expected env in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("env.exe", pathenv)
	if err != nil {
		t.Errorf("expected env.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/env", pathenv)
	if err != nil {
		t.Errorf("expected env in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}

	p, err = lookPath("/usr/bin/env.exe", pathenv)
	if err != nil {
		t.Errorf("expected env.exe in path: %s", pathenv)
	}

	if !strings.Contains(p, "env.exe") {
		t.Errorf("expected env.exe in name: %s", p)
	}
}

func TestCygpathToUnixPathList(t *testing.T) {

	if Cygpath == nil {
		t.Skip(errCygpath.Error())
	}

	// MSYS2: "/c/Windows:/c/Bla"
	// Cywin: "/cygdrive/c/Windows:/cygdrive/c/Bla"
	path := Cygpath.tryUnixPathList("C:\\Windows;C:\\Bla")

	if !strings.Contains(path, "/c/Windows") {
		t.Errorf("PATH: got %v, want %v", path, "/c/Windows")
	}

	if !strings.Contains(path, "/c/Bla") {
		t.Errorf("PATH: got %v, want %v", path, "/c/Bla")
	}

	if !strings.Contains(path, ":") {
		t.Errorf("PATH: got %v, want %v", path, ":")
	}

	tests := []struct {
		inVal  string
		sysVal string
	}{
		{"||C:/Users/user2/AppData/Local/Programs/Git/usr/bin/lesspipe.sh %s", "||C:/Users/user2/AppData/Local/Programs/Git/usr/bin/lesspipe.sh %s"},
		{"-C:\\BlaBla", "-C:\\BlaBla"},
		{"This;Has;No;Paths", "This:Has:No:Paths"},
		{";;", ""},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := Cygpath.tryUnixPathList(tt.inVal); got != tt.sysVal {
				t.Errorf("got %v, want %v", got, tt.sysVal)
			}
		})
	}
}

func BenchmarkCygpathGetEnv(b *testing.B) {
	if Cygpath == nil {
		b.Skip(errCygpath.Error())
	}
	for i := 0; i < b.N; i++ {
		GetEnv()
	}
}
