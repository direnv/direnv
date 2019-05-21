package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	toml "github.com/BurntSushi/toml"
)

type Config struct {
	Env             Env
	WorkDir         string // Current directory
	ConfDir         string
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

func LoadConfig(env Env) (config *Config, err error) {
	config = &Config{
		Env: env,
	}

	config.ConfDir = env[DIRENV_CONFIG]
	if config.ConfDir == "" {
		config.ConfDir = XdgConfigDir(env, "direnv")
	}
	if config.ConfDir == "" {
		err = fmt.Errorf("Couldn't find a configuration directory for direnv")
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

		for _, prefix := range tomlConf.Whitelist.Prefix {
			config.WhitelistPrefix = append(config.WhitelistPrefix, prefix)
		}

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
			log_error("invalid DIRENV_WARN_TIMEOUT: " + err.Error())
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
			err = fmt.Errorf("Can't find bash: %q", err)
			return
		}
	}

	return
}

func (self *Config) AllowDir() string {
	return filepath.Join(self.ConfDir, "allow")
}

func (self *Config) LoadedRC() *RC {
	if self.RCDir == "" {
		log_debug("RCDir is blank - loadedRC is nil")
		return nil
	}
	rcPath := filepath.Join(self.RCDir, ".envrc")

	times_string := self.Env[DIRENV_WATCHES]

	return RCFromEnv(rcPath, times_string, self)
}

func (self *Config) FindRC() *RC {
	return FindRC(self.WorkDir, self)
}

func (self *Config) EnvDiff() (*EnvDiff, error) {
	if self.Env[DIRENV_DIFF] == "" {
		if self.Env[DIRENV_WATCHES] == "" {
			return self.Env.Diff(self.Env), nil
		} else {
			return nil, fmt.Errorf("DIRENV_DIFF is empty")
		}
	}
	return LoadEnvDiff(self.Env[DIRENV_DIFF])
}
