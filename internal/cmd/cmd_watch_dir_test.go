package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchDirSymlinks(t *testing.T) {
	for _, dangling := range []bool{false, true} {
		name := "newer-target"
		if dangling {
			name = "dangling"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "watched")
			if err := os.Mkdir(dir, 0o700); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(root, "target")
			link := filepath.Join(dir, "link")
			if err := os.Symlink(target, link); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			info, err := os.Lstat(link)
			if err != nil {
				t.Fatal(err)
			}
			modified := info.ModTime().Add(2 * time.Second)
			if !dangling {
				if err := os.WriteFile(target, []byte("value"), 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.Chtimes(target, modified, modified); err != nil {
					t.Fatal(err)
				}
			}

			output, err := os.CreateTemp(root, "export")
			if err != nil {
				t.Fatal(err)
			}
			original := os.Stdout
			defer func() { os.Stdout = original }()
			defer func() {
				if err := output.Close(); err != nil {
					t.Error(err)
				}
			}()
			os.Stdout = output
			err = watchDirCommand(Env{}, []string{"watch-dir", "json", dir})
			os.Stdout = original
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(output.Name())
			if err != nil {
				t.Fatal(err)
			}
			var exported map[string]string
			if err := json.Unmarshal(data, &exported); err != nil {
				t.Fatal(err)
			}
			watches := NewFileTimes()
			if err := watches.Unmarshal(exported[DIRENV_WATCHES]); err != nil {
				t.Fatal(err)
			}
			if err := watches.Check(); err != nil {
				t.Fatalf("unchanged watched directory should be current: %v", err)
			}

			if err := os.WriteFile(target, []byte("changed"), 0o600); err != nil {
				t.Fatal(err)
			}
			modified = modified.Add(time.Second)
			if err := os.Chtimes(target, modified, modified); err != nil {
				t.Fatal(err)
			}
			if err := watches.CheckOne(link); err == nil {
				t.Fatal("changed or newly created target should invalidate the watch")
			}
		})
	}
}
