package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	toml "github.com/BurntSushi/toml"
	"github.com/direnv/direnv/xdg"
)

// Config represents the direnv configuration and state.
type Config struct {
	Env             Env
	WorkDir         string // Current directory
	ConfDir         string
	DataDir         string
	SelfPath        string
	BashPath        string
	RCDir           string
	TomlPath        string
	DisableStdin    bool
	WarnTimeout     time.Duration
	WhitelistPrefix []string
	WhitelistExact  map[string]bool
}

type tomlConfig struct {
	BashPath     string        `toml:"bash_path"`
	DisableStdin bool          `toml:"disable_stdin"`
	WarnTimeout  time.Duration `toml:"warn_timeout"`
	Whitelist    whitelist     `toml:"whitelist"`
}

type whitelist struct {
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
		err = fmt.Errorf("LoadConfig() os.Executable() failed: %q", err)
		return
	}
	// Fix for mingsys
	exePath = strings.Replace(exePath, "\\", "/", -1)
	config.SelfPath = exePath

	if config.WorkDir, err = os.Getwd(); err != nil {
		err = fmt.Errorf("LoadConfig() Getwd failed: %q", err)
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
		config.TomlPath = ""
	}

	config.TomlPath = filepath.Join(config.ConfDir, "config.toml")
	if _, statErr := os.Stat(config.TomlPath); statErr != nil {
		config.TomlPath = ""
	}

	if config.TomlPath != "" {
		var tomlConf tomlConfig
		if _, err = toml.DecodeFile(config.TomlPath, &tomlConf); err != nil {
			err = fmt.Errorf("LoadConfig() failed to parse %s: %q", config.TomlPath, err)
			return
		}

		config.WhitelistPrefix = append(config.WhitelistPrefix, tomlConf.Whitelist.Prefix...)

		for _, path := range tomlConf.Whitelist.Exact {
			if !strings.HasSuffix(path, "/.envrc") {
				path = filepath.Join(path, ".envrc")
			}

			config.WhitelistExact[path] = true
		}

		config.DisableStdin = tomlConf.DisableStdin
		config.BashPath = tomlConf.BashPath
		config.WarnTimeout = tomlConf.WarnTimeout
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
			err = fmt.Errorf("can't find bash: %q", err)
			return
		}
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

// FindRC looks for a RC file in the config environment
func (config *Config) FindRC() *RC {
	return FindRC(config.WorkDir, config)
}

// EnvDiff returns the recorded environment diff that was stored if any.
func (config *Config) EnvDiff() (*EnvDiff, error) {
	if config.Env[DIRENV_DIFF] == "" {
		return nil, nil
	}
	return LoadEnvDiff(config.Env[DIRENV_DIFF])
}
