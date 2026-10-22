package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/direnv/direnv/v2/pkg/sri"
)

func sriHashOf(t *testing.T, content string) string {
	t.Helper()

	w := sri.NewWriter(io.Discard, sri.SHA256)
	if _, err := w.Write([]byte(content)); err != nil {
		t.Fatal("could not hash content:", err)
	}
	return w.Sum().String()
}

// captureStdout returns whatever fn writes to os.Stdout.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("could not create pipe:", err)
	}
	os.Stdout = w

	fnErr := fn()

	if err := w.Close(); err != nil {
		t.Fatal("could not close pipe writer:", err)
	}
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal("could not read captured output:", err)
	}
	return string(out), fnErr
}

func fetchURL(t *testing.T, config *Config, url string, hash string) (string, error) {
	t.Helper()

	out, err := captureStdout(t, func() error {
		return cmdFetchURL(Env{}, []string{"fetchurl", url, hash}, config)
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func TestFetchURLChangedURL(t *testing.T) {
	const (
		contentA = "content from URL A\n"
		contentB = "completely different content from URL B\n"
	)
	hashA := sriHashOf(t, contentA)
	hashB := sriHashOf(t, contentB)

	mux := http.NewServeMux()
	mux.HandleFunc("/a.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, contentA)
	})
	mux.HandleFunc("/b.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, contentB)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	urlA := server.URL + "/a.txt"
	urlB := server.URL + "/b.txt"
	config := &Config{CacheDir: t.TempDir()}

	// Fetching A with its own hash downloads and caches it.
	pathA, err := fetchURL(t, config, urlA, hashA)
	if err != nil {
		t.Fatal("fetchurl A failed:", err)
	}
	if data, err := os.ReadFile(pathA); err != nil || string(data) != contentA {
		t.Fatalf("cached file = %q, %v; want %q", data, err, contentA)
	}

	// The same URL and hash again is a cache hit returning the same path.
	pathA2, err := fetchURL(t, config, urlA, hashA)
	if err != nil {
		t.Fatal("fetchurl A (cached) failed:", err)
	}
	if pathA2 != pathA {
		t.Errorf("cache hit returned %q, want %q", pathA2, pathA)
	}

	// Regression test for #1282: keeping A's hash but changing the URL to B
	// must not silently return A's cached content.
	// The new URL misses the cache, the content is downloaded and the
	// integrity check rejects it.
	if _, err := fetchURL(t, config, urlB, hashA); err == nil {
		t.Error("fetchurl B with A's hash succeeded; want a hash mismatch error")
	} else if !strings.Contains(err.Error(), "hash mismatch") {
		t.Errorf("fetchurl B with A's hash: error = %v; want hash mismatch", err)
	}

	// Fetching B with B's own hash works and is a different cache entry.
	pathB, err := fetchURL(t, config, urlB, hashB)
	if err != nil {
		t.Fatal("fetchurl B failed:", err)
	}
	if pathB == pathA {
		t.Error("fetchurl B returned A's cache path")
	}
	if data, err := os.ReadFile(pathB); err != nil || string(data) != contentB {
		t.Fatalf("cached file = %q, %v; want %q", data, err, contentB)
	}
}

func TestCasPath(t *testing.T) {
	hashA, err := sri.Parse(sriHashOf(t, "content A"))
	if err != nil {
		t.Fatal("could not parse hash:", err)
	}
	hashB, err := sri.Parse(sriHashOf(t, "content B"))
	if err != nil {
		t.Fatal("could not parse hash:", err)
	}

	dir := t.TempDir()
	urlA := "https://example.com/a"
	urlB := "https://example.com/b"

	if got, want := casPath(dir, urlA, hashA), casPath(dir, urlA, hashA); got != want {
		t.Errorf("casPath is not stable: %q != %q", got, want)
	}
	if got, other := casPath(dir, urlA, hashA), casPath(dir, urlB, hashA); got == other {
		t.Errorf("casPath does not depend on the URL: both are %q", got)
	}
	if got, other := casPath(dir, urlA, hashA), casPath(dir, urlA, hashB); got == other {
		t.Errorf("casPath does not depend on the hash: both are %q", got)
	}
}
