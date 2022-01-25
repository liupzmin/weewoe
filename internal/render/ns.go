package render

import (
	"errors"
	"strings"

	"github.com/liupzmin/tview"

	"github.com/gdamore/tcell/v2"
	v1 "k8s.io/api/core/v1"
)

// Namespace renders a K8s Namespace to screen.
type Namespace struct {
	Base
}

// ColorerFunc colors a resource row.
func (n Namespace) ColorerFunc() ColorerFunc {
	return func(ns string, h Header, re RowEvent) tcell.Color {
		c := DefaultColorer(ns, h, re)

		if re.Kind == EventUpdate {
			c = StdColor
		}
		if strings.Contains(strings.TrimSpace(re.Row.Fields[0]), "*") {
			c = HighlightColor
		}

		if !Happy(ns, h, re.Row) {
			c = ErrColor
		}

		return c
	}
}

// Header returns a header rbw.
func (Namespace) Header(string) Header {
	return Header{
		HeaderColumn{Name: "NAME"},
		HeaderColumn{Name: "PROCESS", Align: tview.AlignRight},
		HeaderColumn{Name: "HOST", Align: tview.AlignRight},
	}
}

// Render renders a K8s resource to screen.
func (n Namespace) Render(o interface{}, _ string, rows *Rows) error {
	// todo: 待处理
	return nil
}

func (Namespace) diagnose(phase v1.NamespacePhase) error {
	if phase != v1.NamespaceActive && phase != v1.NamespaceTerminating {
		return errors.New("namespace not ready")
	}
	return nil
}
