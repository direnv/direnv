package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/direnv/direnv/v2/pkg/sri"
	"github.com/mattn/go-isatty"
)

// CmdFetchURL is `direnv fetchurl <url> [<integrity-hash>]`
var CmdFetchURL = &Cmd{
	Name:   "fetchurl",
	Desc:   "Fetches a given URL into direnv's CAS",
	Args:   []string{"<url>", "[<integrity-hash>]"},
	Action: actionWithConfig(cmdFetchURL),
}

func cmdFetchURL(_ Env, args []string, config *Config) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("missing URL argument")
	}

	var (
		algo          sri.Algo = sri.SHA256
		url           string
		integrityHash string
	)
	casDir := casDir(config)
	isTTY := isatty.IsTerminal(os.Stdout.Fd())

	url = args[1]
	// Validate the SRI hash if it exists
	if len(args) > 2 {
		// Support Base64 where '/' have been replaced by '_'
		integrityHash = strings.ReplaceAll(args[2], "_", "/")

		hash, err := sri.Parse(integrityHash)
		if err != nil {
			return err
		}

		// Shortcut if the cache already has the file
		casFile := casPath(casDir, hash)
		if fileExists(casFile) {
			fmt.Println(casFile)
			return nil
		}
	}

	// Create the CAS directory if it doesn't exist
	if err = os.MkdirAll(casDir, os.FileMode(0755)); err != nil {
		return err
	}

	// Create a temporary file to copy the content into, before the CAS file
	// location can be calculated.
	tmpfile, err := os.CreateTemp(casDir, "tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()

	// Get the URL
	// G107: Potential HTTP request made with variable url
	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Abort if we don't get a 200 back
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}

	// While copying the content into the temporary location, also calculate the
	// SRI hash.
	w := sri.NewWriter(tmpfile, algo)
	if _, err = io.Copy(w, resp.Body); err != nil {
		return err
	}

	// Here is the new SRI hash
	calculatedHash := w.Sum()

	// Make the file read-only and executable for later
	if err = os.Chmod(tmpfile.Name(), os.FileMode(0500)); err != nil {
		return err
	}

	// Validate if a comparison hash was given
	if integrityHash != "" && calculatedHash.String() != integrityHash {
		return fmt.Errorf("hash mismatch. Expected '%s' but got '%s'", integrityHash, calculatedHash.String())
	}

	// Derive the CAS file location from the SRI hash
	casFile := casPath(casDir, calculatedHash)

	// Put the file into the CAS store if it's not already there
	if !fileExists(casFile) {
		err = tmpfile.Close()
		if err != nil {
			return err
		}
		// Move the temporary file to the CAS location.
		if err = os.Rename(tmpfile.Name(), casFile); err != nil {
			return err
		}
	}

	if integrityHash == "" {
		if isTTY {
			// Print an example for terminal users
			fmt.Printf(`Found hash: %s

Invoke fetchurl again with the hash as an argument to get the disk location:

  direnv fetchurl "%s" "%s"
  #=> %s
`, calculatedHash, url, calculatedHash.String(), casFile)
		} else {
			// Only print the hash in scripting mode. Add one extra hurdle on
			// purpose to use fetchurl without the SRI hash.
			_, err = fmt.Println(calculatedHash)
		}
	} else {
		// Print the location to the CAS file
		_, err = fmt.Println(casFile)
	}
	return err
}

func casDir(c *Config) string {
	return filepath.Join(c.CacheDir, "cas")
}

// casPath returns filesystem path for SRI hashes
func casPath(dir string, integrityHash *sri.Hash) string {
	// Use Hex encoding for the filesystem to avoid issues
	sriFile := integrityHash.Hex()
	return filepath.Join(dir, sriFile)
}
