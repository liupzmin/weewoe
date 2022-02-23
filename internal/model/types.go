package model

import (
	"context"

	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/sahilm/fuzzy"
)

type Cluster struct {
	Total, Healthy, Killed, PortClosed, Node int
}

// ResourceViewerListener listens to viewing resource events.
type ResourceViewerListener interface {
	ResourceChanged(lines []string, matches fuzzy.Matches)
	ResourceFailed(error)
}

// ViewerToggleOpts represents a collection of viewing options.
type ViewerToggleOpts map[string]bool

// ResourceViewer represents a viewed resource.
type ResourceViewer interface {
	GetPath() string
	Filter(string)
	ClearFilter()
	Peek() []string
	SetOptions(context.Context, ViewerToggleOpts)
	Watch(context.Context) error
	Refresh(context.Context) error
	AddListener(ResourceViewerListener)
	RemoveListener(ResourceViewerListener)
}

// Igniter represents a runnable view.
type Igniter interface {
	Init(ctx context.Context) error

	// Start starts a component.
	Start()

	// Stop terminates a component.
	Stop()
}

// Hinter represent a menu mnemonic provider.
type Hinter interface {
	// Hints returns a collection of menu hints.
	Hints() MenuHints

	// ExtraHints returns additional hints.
	ExtraHints() map[string]string
}

// Primitive represents a UI primitive.
type Primitive interface {
	tview.Primitive

	// Name returns the view name.
	Name() string
}

// Commander tracks prompt status.
type Commander interface {
	// InCmdMode checks if prompt is active.
	InCmdMode() bool
}

// Component represents a ui component.
type Component interface {
	Primitive
	Igniter
	Hinter
	Commander
}

// Renderer represents a resource renderer.
type Renderer interface {
	// IsGeneric identifies a generic handler.
	IsGeneric() bool

	// Render converts raw resources to tabular data.
	Render(o interface{}, ns string, rows *render.Rows) error

	// Header returns the resource header.
	Header(ns string) render.Header

	// ColorerFunc returns a row colorer function.
	ColorerFunc() render.ColorerFunc
}

// TreeRenderer represents an xray node.
type TreeRenderer interface {
	Render(ctx context.Context, ns string, o interface{}) error
}

// ResourceMeta represents model info about a resource.
type ResourceMeta struct {
	//DAO          dao.Accessor
	Renderer     Renderer
	TreeRenderer TreeRenderer
}
