package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Core        CoreConfig        `toml:"core"`
	Hal         HalConfig         `toml:"hal"`
	Integration IntegrationConfig `toml:"integration"`
}

type CoreConfig struct {
	Editor    string `toml:"editor"`
	SelectCmd string `toml:"selectcmd"`
	TomlFile  string `toml:"tomlfile"`
}

type HalConfig struct {
	Token   string `toml:"token"`
	BaseURL string `toml:"base_url"`
	Dir     string `toml:"dir"`
	Repo    string `toml:"repo"`
}

type IntegrationConfig struct {
	GitHubToken  string `toml:"github_token"`
	SlackToken   string `toml:"slack_token"`
	DataDogToken string `toml:"datadog_token"`
}

var Conf Config

func GetDefaultDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	}
	dir = filepath.Join(dir, "hal-ops")

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return dir, fmt.Errorf("cannot create directory: %v", err)
	}

	return dir, nil
}

func (cfg *Config) LoadFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	cfg.Core.Editor = os.Getenv("EDITOR")
	if cfg.Core.Editor == "" {
		cfg.Core.Editor = "vim"
	}
	cfg.Core.SelectCmd = "fzf-tmux --multi:fzf --multi:peco"
	cfg.Core.TomlFile = file

	cfg.Hal.Token = os.Getenv("GITHUB_TOKEN")
	cfg.Hal.BaseURL = "https://hal-ops.github.com"
	dir := filepath.Join(filepath.Dir(file), "files")
	os.MkdirAll(dir, 0700)
	cfg.Hal.Dir = dir

	return toml.NewEncoder(f).Encode(cfg)
}
