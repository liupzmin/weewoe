package model

import (
	"github.com/liupzmin/weewoe/internal/render"
)

// TableListener represents a table model listener.
type TableListener interface {
	// TableDataChanged notifies the model data changed.
	TableDataChanged(render.TableData)

	// TableLoadFailed notifies the load failed.
	TableLoadFailed(error)
}
