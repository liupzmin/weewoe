package config_test

import (
	"testing"

	"github.com/liupzmin/weewoe/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIsReadOnly(t *testing.T) {
	uu := map[string]struct {
		config      string
		read, write bool
		readOnly    bool
	}{
		"writable": {
			config: "w2.yml",
		},
		"writable_read_override": {
			config:   "w2.yml",
			read:     true,
			readOnly: true,
		},
		"writable_write_override": {
			config: "w2.yml",
			write:  true,
		},
		"readonly": {
			config:   "w2_readonly.yml",
			readOnly: true,
		},
		"readonly_read_override": {
			config:   "w2_readonly.yml",
			read:     true,
			readOnly: true,
		},
		"readonly_write_override": {
			config: "w2_readonly.yml",
			write:  true,
		},
		"readonly_both_override": {
			config: "w2_readonly.yml",
			read:   true,
			write:  true,
		},
	}

	cfg := config.NewConfig()
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Nil(t, cfg.Load("testdata/"+u.config))
			cfg.W2.OverrideReadOnly(u.read)
			cfg.W2.OverrideWrite(u.write)
			assert.Equal(t, u.readOnly, cfg.W2.IsReadOnly())
		})
	}
}

func TestK9sActiveClusterBlank(t *testing.T) {
	var c config.W2
	cl := c.ActiveCluster()
	assert.Equal(t, config.NewCluster(), cl)
}

func TestK9sActiveCluster(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, cfg.Load("testdata/k9s.yml"))

	cl := cfg.W2.ActiveCluster()
	assert.NotNil(t, cl)
	assert.Equal(t, "kube-system", cl.Namespace.Active)
	assert.Equal(t, 5, len(cl.Namespace.Favorites))
}

func TestGetScreenDumpDir(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, cfg.Load("testdata/k9s.yml"))

	assert.Equal(t, "/tmp", cfg.W2.GetScreenDumpDir())
}

func TestGetScreenDumpDirOverride(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, cfg.Load("testdata/k9s.yml"))
	cfg.W2.OverrideScreenDumpDir("/override")

	assert.Equal(t, "/override", cfg.W2.GetScreenDumpDir())
}

func TestGetScreenDumpDirOverrideEmpty(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, cfg.Load("testdata/k9s.yml"))
	cfg.W2.OverrideScreenDumpDir("")

	assert.Equal(t, "/tmp", cfg.W2.GetScreenDumpDir())
}

func TestGetScreenDumpDirEmpty(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, cfg.Load("testdata/k9s1.yml"))
	cfg.W2.OverrideScreenDumpDir("")

	assert.Equal(t, config.W2DefaultScreenDumpDir, cfg.W2.GetScreenDumpDir())
}
