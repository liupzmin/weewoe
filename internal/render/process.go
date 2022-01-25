package render

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Process renders a K8s Pod to screen.
type Process struct {
	Base
}

// ColorerFunc colors a resource row.
func (p Process) ColorerFunc() ColorerFunc {
	return func(ns string, h Header, re RowEvent) tcell.Color {
		c := DefaultColorer(ns, h, re)

		statusCol := h.IndexOf("STATUS", true)
		pstatusCol := h.IndexOf("PORTS-STATUS", true)

		if statusCol == -1 && pstatusCol == -1 {
			return c
		}
		status := strings.TrimSpace(re.Row.Fields[statusCol])
		switch status {
		case Pending:
			c = PendingColor
		case ContainerCreating, PodInitializing:
			c = AddColor
		case Initialized:
			c = HighlightColor
		case Completed:
			c = CompletedColor
		case Running:
			c = StdColor
			if !Happy(ns, h, re.Row) {
				c = ErrColor
			}
		case KILLED:
			c = ErrColor
		case OPEN:
			c = StdColor
		case CLOSED:
			c = KillColor
		default:
			if !Happy(ns, h, re.Row) {
				c = ErrColor
			}
		}
		return c
	}
}

// Header returns a header row.
func (Process) Header(ns string) Header {
	return Header{
		HeaderColumn{Name: "NAMESPACE"},
		HeaderColumn{Name: "NAME"},
		HeaderColumn{Name: "STATUS"},
		HeaderColumn{Name: "NODE"},
		HeaderColumn{Name: "PORTS"},
		HeaderColumn{Name: "PORTS-STATUS"},
		HeaderColumn{Name: "PATH"},
		HeaderColumn{Name: "FLAG"},
		HeaderColumn{Name: "START-TIME"},
		HeaderColumn{Name: "UPDATE-TIME"},
		HeaderColumn{Name: "Suspend"},
	}
}

// Render renders a K8s resource to screen.
func (p Process) Render(o interface{}, ns string, rows *Rows) error {
	// todo: 待处理
	return nil
}

func (p Process) diagnose(phase string, cr, ct int) error {
	if phase == Completed {
		return nil
	}
	if cr != ct || ct == 0 {
		return fmt.Errorf("container ready check failed: %d of %d", cr, ct)
	}

	return nil
}
