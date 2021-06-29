package main

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
	BashBuiltin     bool
	BashPath        string
	RCDir           string
	TomlPath        string
	DisableStdin    bool
	StrictEnv       bool
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
	BashBuiltin  bool         `toml:"bash_builtin"`
	BashPath     string       `toml:"bash_path"`
	DisableStdin bool         `toml:"disable_stdin"`
	StrictEnv    bool         `toml:"strict_env"`
	WarnTimeout  tomlDuration `toml:"warn_timeout"`
}

type tomlWhitelist struct {
	Prefix []string
	Exact  []string
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

	config.RCDir = env[DIRENV_DIR]
	if len(config.RCDir) > 0 && config.RCDir[0:1] == "-" {
		config.RCDir = config.RCDir[1:]
	}

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
		var global tomlGlobal
		tomlConf := tomlConfig{
			tomlGlobal: &global,
			Global:     &global,
		}
		if _, err = toml.DecodeFile(config.TomlPath, &tomlConf); err != nil {
			err = fmt.Errorf("LoadConfig() failed to parse %s: %w", config.TomlPath, err)
			return
		}

		config.WhitelistPrefix = append(config.WhitelistPrefix, tomlConf.Whitelist.Prefix...)

		for _, path := range tomlConf.Whitelist.Exact {
			if !strings.HasSuffix(path, "/.envrc") {
				path = filepath.Join(path, ".envrc")
			}

			config.WhitelistExact[path] = true
		}

		config.BashBuiltin = tomlConf.BashBuiltin
		config.BashPath = tomlConf.BashPath
		config.DisableStdin = tomlConf.DisableStdin
		config.StrictEnv = tomlConf.StrictEnv
		config.WarnTimeout = tomlConf.WarnTimeout.Duration
	}

	if config.WarnTimeout == 0 {
		timeout, err := time.ParseDuration(env.Fetch("DIRENV_WARN_TIMEOUT", "5s"))
		if err != nil {
			logError("invalid DIRENV_WARN_TIMEOUT: " + err.Error())
			timeout = 5 * time.Second
		}
		config.WarnTimeout = timeout
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

// LoadedRC returns a RC file if any has been loaded
func (config *Config) LoadedRC() *RC {
	if config.RCDir == "" {
		logDebug("RCDir is blank - loadedRC is nil")
		return nil
	}
	rcPath := filepath.Join(config.RCDir, ".envrc")

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
