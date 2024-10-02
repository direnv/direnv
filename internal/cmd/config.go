package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	toml "github.com/BurntSushi/toml"
	"github.com/direnv/direnv/v2/xdg"
)

// Config represents the direnv configuration and state.
type Config struct {
	Env             Env
	WorkDir         string // Current directory
	ConfDir         string
	CacheDir        string
	DataDir         string
	SelfPath        string
	BashPath        string
	RCFile          string
	TomlPath        string
	HideEnvDiff     bool
	DisableStdin    bool
	StrictEnv       bool
	LoadDotenv      bool
	WarnTimeout     time.Duration
	WhitelistPrefix []string
	WhitelistExact  map[string]bool
}

type tomlDuration struct {
	time.Duration
}

func (d *tomlDuration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type tomlConfig struct {
	*tomlGlobal               // For backward-compatibility
	Global      *tomlGlobal   `toml:"global"`
	Whitelist   tomlWhitelist `toml:"whitelist"`
}

type tomlGlobal struct {
	BashPath     string        `toml:"bash_path"`
	DisableStdin bool          `toml:"disable_stdin"`
	StrictEnv    bool          `toml:"strict_env"`
	SkipDotenv   bool          `toml:"skip_dotenv"` // deprecated, use load_dotenv
	LoadDotenv   bool          `toml:"load_dotenv"`
	WarnTimeout  *tomlDuration `toml:"warn_timeout"`
	HideEnvDiff  bool          `toml:"hide_env_diff"`
}

type tomlWhitelist struct {
	Prefix []string `toml:"prefix"`
	Exact  []string `toml:"exact"`
}

// Expand a path string prefixed with ~/ to the current user's home directory.
// Example: if current user is user1 with home directory in /home/user1, then
// ~/project -> /home/user1/project
// It's useful to allow paths with ~/, so that direnv.toml can be reused via
// dotfiles repos across systems with different standard home paths
// (compare Linux /home and macOS /Users).
func expandTildePath(path string) (pathExpanded string) {
	pathExpanded = path
	if strings.HasPrefix(path, "~/") {
		if homedir, homedirErr := os.UserHomeDir(); homedirErr == nil {
			pathExpanded = filepath.Join(homedir, path[2:])
		}
	}
	return pathExpanded
}

// LoadConfig opens up the direnv configuration from the Env.
func LoadConfig(env Env) (config *Config, err error) {
	config = &Config{
		Env: env,
	}

	config.ConfDir = env[DIRENV_CONFIG]
	if config.ConfDir == "" {
		config.ConfDir = xdg.ConfigDir(env, "direnv")
	}
	if config.ConfDir == "" {
		err = fmt.Errorf("couldn't find a configuration directory for direnv")
		return
	}

	var exePath string
	if exePath, err = os.Executable(); err != nil {
		err = fmt.Errorf("LoadConfig() os.Executable() failed: %w", err)
		return
	}
	// Fix for mingsys
	exePath = strings.Replace(exePath, "\\", "/", -1)
	config.SelfPath = exePath

	if config.WorkDir, err = os.Getwd(); err != nil {
		err = fmt.Errorf("LoadConfig() Getwd failed: %w", err)
		return
	}

	// Default Warn Timeout
	config.WarnTimeout = 5 * time.Second

	config.RCFile = env[DIRENV_FILE]

	config.WhitelistPrefix = make([]string, 0)
	config.WhitelistExact = make(map[string]bool)

	// Load the TOML config
	config.TomlPath = filepath.Join(config.ConfDir, "direnv.toml")
	if _, statErr := os.Stat(config.TomlPath); statErr != nil {
		config.TomlPath = filepath.Join(config.ConfDir, "config.toml")
		if _, statErr := os.Stat(config.TomlPath); statErr != nil {
			config.TomlPath = ""
		}
	}

	if config.TomlPath != "" {
		// Declare global once and then share it between the top-level and Global
		// keys. The goal here is to let the decoder fill global regardless of if
		// the values are in the [global] section or not. The reason we do that is
		// to keep backward-compatibility with the old top-level notation.
		var global tomlGlobal
		tomlConf := tomlConfig{
			tomlGlobal: &global,
			Global:     &global,
		}
		if _, err = toml.DecodeFile(config.TomlPath, &tomlConf); err != nil {
			err = fmt.Errorf("LoadConfig() failed to parse %s: %w", config.TomlPath, err)
			return
		}

		config.HideEnvDiff = tomlConf.HideEnvDiff

		for _, path := range tomlConf.Whitelist.Prefix {
			config.WhitelistPrefix = append(config.WhitelistPrefix, expandTildePath(path))
		}

		for _, path := range tomlConf.Whitelist.Exact {
			if !(strings.HasSuffix(path, "/.envrc") || strings.HasSuffix(path, "/.env")) {
				path = filepath.Join(path, ".envrc")
			}

			config.WhitelistExact[expandTildePath(path)] = true
		}

		if tomlConf.SkipDotenv {
			logError("skip_dotenv has been inverted to load_dotenv.")
		}

		config.BashPath = tomlConf.BashPath
		config.DisableStdin = tomlConf.DisableStdin
		config.LoadDotenv = tomlConf.LoadDotenv
		config.StrictEnv = tomlConf.StrictEnv
		if tomlConf.WarnTimeout != nil {
			config.WarnTimeout = tomlConf.WarnTimeout.Duration
		}
	}

	if ts := env.Fetch("DIRENV_WARN_TIMEOUT", ""); ts != "" {
		timeout, err := time.ParseDuration(ts)
		if err == nil {
			config.WarnTimeout = timeout
		} else {
			logError("invalid DIRENV_WARN_TIMEOUT: " + err.Error())
		}
	}

	if config.BashPath == "" {
		if env[DIRENV_BASH] != "" {
			config.BashPath = env[DIRENV_BASH]
		} else if bashPath != "" {
			config.BashPath = bashPath
		} else if config.BashPath, err = exec.LookPath("bash"); err != nil {
			err = fmt.Errorf("can't find bash: %w", err)
			return
		}
	}

	if config.CacheDir == "" {
		config.CacheDir = xdg.CacheDir(env, "direnv")
	}
	if config.CacheDir == "" {
		err = fmt.Errorf("couldn't find a cache directory for direnv")
		return
	}

	if config.DataDir == "" {
		config.DataDir = xdg.DataDir(env, "direnv")
	}
	if config.DataDir == "" {
		err = fmt.Errorf("couldn't find a data directory for direnv")
		return
	}

	return
}

// AllowDir is the folder where all the "allow" files are stored.
func (config *Config) AllowDir() string {
	return filepath.Join(config.DataDir, "allow")
}

// DenyDir is the folder where all the "deny" files are stored.
func (config *Config) DenyDir() string {
	return filepath.Join(config.DataDir, "deny")
}

// LoadedRC returns a RC file if any has been loaded
func (config *Config) LoadedRC() *RC {
	if config.Env[DIRENV_FILE] == "" {
		logDebug("RCFile is blank - loadedRC is nil")
		return nil
	}
	rcPath := config.Env[DIRENV_FILE]

	timesString := config.Env[DIRENV_WATCHES]

	return RCFromEnv(rcPath, timesString, config)
}

// EnvFromRC loads an RC from a specified path and returns the new environment
func (config *Config) EnvFromRC(path string, previousEnv Env) (Env, error) {
	rc, err := RCFromPath(path, config)
	if err != nil {
		return nil, err
	}
	return rc.Load(previousEnv)
}

// FindRC looks for a RC file in the config environment
func (config *Config) FindRC() (*RC, error) {
	return FindRC(config.WorkDir, config)
}

// Revert undoes the recorded changes (if any) to the supplied environment,
// returning a new environment
func (config *Config) Revert(env Env) (Env, error) {
	if config.Env[DIRENV_DIFF] == "" {
		return env.Copy(), nil
	}
	diff, err := LoadEnvDiff(config.Env[DIRENV_DIFF])
	if err == nil {
		return diff.Reverse().Patch(env), nil
	}
	return nil, err
}
