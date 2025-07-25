// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"context"
	"io"
	"time"

	"github.com/liupzmin/weewoe/internal/render"
)

type FactoryFn func() MyFactory

type MyFactory interface {
	Stream(cat string) <-chan render.Rows
	SendCommand(int64) error
	Terminate()
}

//// Getter represents a resource getter.
//type Getter interface {
//	// Get return a given resource.
//	Get(ctx context.Context, path string) (runtime.Object, error)
//}
//
//// Lister represents a resource lister.
//type Lister interface {
//	// List returns a resource collection.
//	List(ctx context.Context, ns string) ([]runtime.Object, error)
//}

//// Accessor represents an accessible k8s resource.
//type Accessor interface {
//	Lister
//	Getter
//
//	// Init the resource with a factory object.
//	Init(MyFactory, string)
//
//	Cat() string
//}

// DrainOptions tracks drain attributes.
type DrainOptions struct {
	GracePeriodSeconds  int
	Timeout             time.Duration
	IgnoreAllDaemonSets bool
	DeleteEmptyDirData  bool
	Force               bool
}

// NodeMaintainer performs node maintenance operations.
type NodeMaintainer interface {
	// ToggleCordon toggles cordon/uncordon a node.
	ToggleCordon(path string, cordon bool) error

	// Drain drains the given node.
	Drain(path string, opts DrainOptions, w io.Writer) error
}

// Describer describes a resource.
type Describer interface {
	// Describe describes a resource.
	Describe(path string) (string, error)

	// ToYAML dumps a resource to YAML.
	ToYAML(path string, showManaged bool) (string, error)
}

// Scalable represents resources that can scale.
type Scalable interface {
	// Scale scales a resource up or down.
	Scale(ctx context.Context, path string, replicas int32) error
}

// Controller represents a pod controller.
type Controller interface {
	// Pod returns a pod instance matching the selector.
	Pod(path string) (string, error)
}

// Nuker represents a resource deleter.
type Nuker interface {
	// Delete removes a resource from the api server.
	Delete(path string, cascade, force bool) error
}

// Switchable represents a switchable resource.
type Switchable interface {
	// Switch changes the active context.
	Switch(ctx string) error
}

// Restartable represents a restartable resource.
type Restartable interface {
	// Restart performs a rollout restart.
	Restart(ctx context.Context, path string) error
}

// Runnable represents a runnable resource.
type Runnable interface {
	// Run triggers a run.
	Run(path string) error
}

//// Logger represents a resource that exposes logs.
//type Logger interface {
//	// Logs tails a resource logs.
//	Logs(path string, opts *v1.PodLogOptions) (*restclient.Request, error)
//}
