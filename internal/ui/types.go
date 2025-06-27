// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"context"
	"time"

	"github.com/liupzmin/weewoe/internal/model"
	"github.com/liupzmin/weewoe/internal/render"
)

type (
	// SortFn represent a function that can sort columnar data.
	SortFn func(rows render.Rows, sortCol SortColumn)

	// SortColumn represents a sortable column.
	SortColumn struct {
		name string
		asc  bool
	}
)

// Namespaceable represents a namespaceable model.
type Namespaceable interface {
	// ClusterWide returns true if the model represents resource in all namespaces.
	ClusterWide() bool

	// GetNamespace returns the model namespace.
	GetNamespace() string

	// SetNamespace changes the model namespace.
	SetNamespace(string)

	// InNamespace check if current namespace matches models.
	InNamespace(string) bool
}

// Tabular represents a tabular model.
type Tabular interface {
	Namespaceable

	// SetInstance sets parent resource path.
	SetInstance(string)

	// SetLabelFilter sets the label filter.
	SetLabelFilter(string)

	// Empty returns true if model has no data.
	Empty() bool

	// Count returns the model data count.
	Count() int

	// Peek returns current model data.
	Peek() render.TableData

	// Watch watches a given resource for changes.
	Watch(context.Context) error

	// Refresh forces a new refresh.
	Refresh(context.Context) error

	// SendCommand customs actions
	SendCommand(context.Context, int64) error

	// SetRefreshRate sets the model watch loop rate.
	SetRefreshRate(time.Duration)

	// AddListener registers a model listener.
	AddListener(model.TableListener)

	// RemoveListener unregister a model listener.
	RemoveListener(model.TableListener)

	// Delete a resource.
	Delete(ctx context.Context, path string, cascade, force bool) error
}

type TableListener interface {
	TableTick()
}
