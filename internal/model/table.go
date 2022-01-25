package model

import (
	"fmt"
	"time"

	"github.com/liupzmin/weewoe/internal/render"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
)

const initRefreshRate = 300 * time.Millisecond

// TableListener represents a table model listener.
type TableListener interface {
	// TableDataChanged notifies the model data changed.
	TableDataChanged(render.TableData)

	// TableLoadFailed notifies the load failed.
	TableLoadFailed(error)
}

// ----------------------------------------------------------------------------
// Helpers...

func hydrate(ns string, oo interface{}, rr *render.Rows, re Renderer) error {

	if err := re.Render(oo, ns, rr); err != nil {
		return err
	}

	return nil
}

type Generic interface {
	SetTable(*metav1beta1.Table)
	Header(string) render.Header
	Render(interface{}, string, *render.Row) error
}

func genericHydrate(ns string, table *metav1beta1.Table, rr render.Rows, re Renderer) error {
	gr, ok := re.(Generic)
	if !ok {
		return fmt.Errorf("expecting generic renderer but got %T", re)
	}
	gr.SetTable(table)
	for i, row := range table.Rows {
		if err := gr.Render(row, ns, &rr[i]); err != nil {
			return err
		}
	}

	return nil
}
