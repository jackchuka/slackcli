package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Workspace struct {
	Name   string `json:"name"`
	Token  string `json:"token,omitempty"`
	TeamID string `json:"team_id,omitempty"`
}

type Config struct {
	ActiveWorkspace string               `json:"active_workspace"`
	Workspaces      map[string]Workspace `json:"workspaces"`
	path            string
}

func DefaultPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "slackcli", "config.json")
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath()
	}
	cfg := &Config{
		Workspaces: make(map[string]Workspace),
		path:       path,
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	cfg.path = path
	if cfg.Workspaces == nil {
		cfg.Workspaces = make(map[string]Workspace)
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}

func (c *Config) ActiveToken() string {
	if ws, ok := c.Workspaces[c.ActiveWorkspace]; ok {
		return ws.Token
	}
	return ""
}

func (c *Config) SetWorkspace(name string, ws Workspace) {
	c.Workspaces[name] = ws
	c.ActiveWorkspace = name
}

func (c *Config) RemoveWorkspace(name string) {
	delete(c.Workspaces, name)
	if c.ActiveWorkspace == name {
		c.ActiveWorkspace = ""
		for k := range c.Workspaces {
			c.ActiveWorkspace = k
			break
		}
	}
}
