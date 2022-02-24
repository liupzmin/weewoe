package config

import "fmt"

const (
	defaultRefreshRate  = 2
	defaultMaxConnRetry = 5
)

// W2 tracks W2 configuration options.
type W2 struct {
	RefreshRate         int      `yaml:"refreshRate"`
	MaxConnRetry        int      `yaml:"maxConnRetry"`
	EnableMouse         bool     `yaml:"enableMouse"`
	Headless            bool     `yaml:"headless"`
	Logoless            bool     `yaml:"logoless"`
	Crumbsless          bool     `yaml:"crumbsless"`
	ReadOnly            bool     `yaml:"readOnly"`
	NoIcons             bool     `yaml:"noIcons"`
	Logger              *Logger  `yaml:"logger"`
	Cluster             *Cluster `yaml:"cluster,omitempty"`
	ScreenDumpDir       string   `yaml:"screenDumpDir"`
	host                string
	port                string
	manualRefreshRate   int
	manualHeadless      *bool
	manualLogoless      *bool
	manualCrumbsless    *bool
	manualReadOnly      *bool
	manualCommand       *string
	manualScreenDumpDir *string
}

// NewWeeWoe create a new W2 configuration.
func NewWeeWoe() *W2 {
	return &W2{
		RefreshRate:   defaultRefreshRate,
		MaxConnRetry:  defaultMaxConnRetry,
		Logger:        NewLogger(),
		ScreenDumpDir: W2DefaultScreenDumpDir,
	}
}

// OverrideHost set the host manually.
func (w *W2) OverrideHost(h string) {
	w.host = h
}

// OverridePort set the port manually.
func (w *W2) OverridePort(p string) {
	w.port = p
}

// OverrideRefreshRate set the refresh rate manually.
func (w *W2) OverrideRefreshRate(r int) {
	w.manualRefreshRate = r
}

// OverrideHeadless toggle the header manually.
func (w *W2) OverrideHeadless(b bool) {
	w.manualHeadless = &b
}

// OverrideLogoless toggle the k9s logo manually.
func (w *W2) OverrideLogoless(b bool) {
	w.manualLogoless = &b
}

// OverrideCrumbsless tooh the crumbslessness manually.
func (w *W2) OverrideCrumbsless(b bool) {
	w.manualCrumbsless = &b
}

// OverrideReadOnly set the readonly mode manually.
func (w *W2) OverrideReadOnly(b bool) {
	if b {
		w.manualReadOnly = &b
	}
}

// OverrideWrite set the write mode manually.
func (w *W2) OverrideWrite(b bool) {
	if b {
		var flag bool
		w.manualReadOnly = &flag
	}
}

// OverrideCommand set the command manually.
func (w *W2) OverrideCommand(cmd string) {
	w.manualCommand = &cmd
}

// OverrideScreenDumpDir set the screen dump dir manually.
func (w *W2) OverrideScreenDumpDir(dir string) {
	w.manualScreenDumpDir = &dir
}

// IsHeadless returns headless setting.
func (w *W2) IsHeadless() bool {
	h := w.Headless
	if w.manualHeadless != nil && *w.manualHeadless {
		h = *w.manualHeadless
	}

	return h
}

// IsLogoless returns logoless setting.
func (w *W2) IsLogoless() bool {
	h := w.Logoless
	if w.manualLogoless != nil && *w.manualLogoless {
		h = *w.manualLogoless
	}

	return h
}

// IsCrumbsless returns crumbsless setting.
func (w *W2) IsCrumbsless() bool {
	h := w.Crumbsless
	if w.manualCrumbsless != nil && *w.manualCrumbsless {
		h = *w.manualCrumbsless
	}

	return h
}

// GetRefreshRate returns the current refresh rate.
func (w *W2) GetRefreshRate() int {
	rate := w.RefreshRate
	if w.manualRefreshRate != 0 {
		rate = w.manualRefreshRate
	}

	return rate
}

// IsReadOnly returns the readonly setting.
func (w *W2) IsReadOnly() bool {
	readOnly := w.ReadOnly
	if w.manualReadOnly != nil {
		readOnly = *w.manualReadOnly
	}

	return readOnly
}

// ActiveCluster returns the currently active cluster.
func (w *W2) ActiveCluster() *Cluster {
	if w.Cluster == nil {
		w.Cluster = NewCluster()
	}

	return w.Cluster
}

func (w *W2) GetScreenDumpDir() string {
	screenDumpDir := w.ScreenDumpDir

	if w.manualScreenDumpDir != nil && *w.manualScreenDumpDir != "" {
		screenDumpDir = *w.manualScreenDumpDir
	}

	if screenDumpDir == "" {
		return W2DefaultScreenDumpDir
	}

	return screenDumpDir
}

func (w *W2) validateDefaults() {
	if w.RefreshRate <= 0 {
		w.RefreshRate = defaultRefreshRate
	}
	if w.MaxConnRetry <= 0 {
		w.MaxConnRetry = defaultMaxConnRetry
	}
	if w.ScreenDumpDir == "" {
		w.ScreenDumpDir = W2DefaultScreenDumpDir
	}
}

func (w *W2) GRPCServer() string {
	return fmt.Sprintf("%s:%s", w.host, w.port)
}
