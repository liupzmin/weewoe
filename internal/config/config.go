package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/liupzmin/weewoe/internal/client"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// W2Config represents W2 configuration dir env var.
const W2Config = "W2CONFIG"

var (
	// W2ConfigFile represents W2 config file location.
	W2ConfigFile = filepath.Join(W2Home(), "config.yml")
	// W2DefaultScreenDumpDir represents a default directory where W2 screen dumps will be persisted.
	W2DefaultScreenDumpDir = filepath.Join(os.TempDir(), fmt.Sprintf("weewoe-screens-%s", MustW2User()))
)

type (
	// Config tracks W2 configuration options.
	Config struct {
		W2 *W2 `yaml:"weewoe"`
	}
)

// W2Home returns k9s configs home directory.
func W2Home() string {
	if env := os.Getenv(W2Config); env != "" {
		return env
	}

	xdgW2Home, err := xdg.ConfigFile("w2")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create configuration directory for k9s")
	}

	return xdgW2Home
}

// NewConfig creates a new default config.
func NewConfig() *Config {
	return &Config{W2: NewWeeWoe()}
}

// Refine the configuration based on cli args.
func (c *Config) Refine() error {

	var ns = client.DefaultNamespace

	ns = c.W2.ActiveCluster().Namespace.Active

	if err := c.SetActiveNamespace(ns); err != nil {
		return err
	}

	EnsurePath(c.W2.GetScreenDumpDir(), DefaultDirMod)

	return nil
}

// ActiveNamespace returns the active namespace in the current cluster.
func (c *Config) ActiveNamespace() string {
	if c.W2.Cluster == nil {
		log.Warn().Msgf("No context detected returning default namespace")
		return "default"
	}
	cl := c.W2.ActiveCluster()
	if cl != nil && cl.Namespace != nil {
		return cl.Namespace.Active
	}
	if cl == nil {
		cl = NewCluster()
		c.W2.Cluster = cl
	}

	return "default"
}

// FavNamespaces returns fav namespaces in the current cluster.
func (c *Config) FavNamespaces() []string {
	cl := c.W2.ActiveCluster()

	return cl.Namespace.Favorites
}

// SetActiveNamespace set the active namespace in the current cluster.
func (c *Config) SetActiveNamespace(ns string) error {
	if cl := c.W2.ActiveCluster(); cl != nil {
		return cl.Namespace.SetActive(ns)
	}
	err := errors.New("no active cluster. unable to set active namespace")
	log.Error().Err(err).Msg("SetActiveNamespace")

	return err
}

// ActiveView returns the active view in the current cluster.
func (c *Config) ActiveView() string {
	cl := c.W2.ActiveCluster()
	if cl == nil {
		return defaultView
	}
	cmd := cl.View.Active
	if c.W2.manualCommand != nil && *c.W2.manualCommand != "" {
		cmd = *c.W2.manualCommand
	}

	return cmd
}

// SetActiveView set the currently cluster active view.
func (c *Config) SetActiveView(view string) {
	if cl := c.W2.ActiveCluster(); cl != nil {
		cl.View.Active = view
	}
}

// Load W2 configuration from file.
func (c *Config) Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	c.W2 = NewWeeWoe()

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return err
	}
	if cfg.W2 != nil {
		c.W2 = cfg.W2
	}
	if c.W2.Logger == nil {
		c.W2.Logger = NewLogger()
	}
	return nil
}

// Save configuration to disk.
func (c *Config) Save() error {
	return c.SaveFile(W2ConfigFile)
}

// SaveFile W2 configuration to disk.
func (c *Config) SaveFile(path string) error {
	EnsurePath(path, DefaultDirMod)
	cfg, err := yaml.Marshal(c)
	if err != nil {
		log.Error().Msgf("[Config] Unable to save W2 config file: %v", err)
		return err
	}
	return os.WriteFile(path, cfg, 0644)
}

// Dump debug...
func (c *Config) Dump(msg string) {
	log.Debug().Msgf("W2 cluster: %s\n", c.W2.Cluster.Namespace)
}

// ----------------------------------------------------------------------------
// Helpers...

func isSet(s *string) bool {
	return s != nil && len(*s) > 0
}
