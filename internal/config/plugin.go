// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// K9sPlugins manages W2 plugins.
var K9sPlugins = filepath.Join(W2Home(), "plugin.yml")

// Plugins represents a collection of plugins.
type Plugins struct {
	Plugin map[string]Plugin `yaml:"plugin"`
}

// Plugin describes a W2 plugin.
type Plugin struct {
	Scopes      []string `yaml:"scopes"`
	Args        []string `yaml:"args"`
	ShortCut    string   `yaml:"shortCut"`
	Pipes       []string `yaml:"pipes"`
	Description string   `yaml:"description"`
	Command     string   `yaml:"command"`
	Confirm     bool     `yaml:"confirm"`
	Background  bool     `yaml:"background"`
}

func (p Plugin) String() string {
	return fmt.Sprintf("[%s] %s(%s)", p.ShortCut, p.Command, strings.Join(p.Args, " "))
}

// NewPlugins returns a new plugin.
func NewPlugins() Plugins {
	return Plugins{
		Plugin: make(map[string]Plugin),
	}
}

// Load W2 plugins.
func (p Plugins) Load() error {
	return p.LoadPlugins(K9sPlugins)
}

// LoadPlugins loads plugins from a given file.
func (p Plugins) LoadPlugins(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var pp Plugins
	if err := yaml.Unmarshal(f, &pp); err != nil {
		return err
	}
	for k, v := range pp.Plugin {
		p.Plugin[k] = v
	}

	return nil
}
