package view

import (
	"time"

	"github.com/liupzmin/weewoe/internal/dao"

	"github.com/liupzmin/weewoe/internal/watch"

	"google.golang.org/grpc"

	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/rs/zerolog/log"
)

const (
	favNSIndicator     = "+"
	defaultNSIndicator = "(*)"
)

// Namespace represents a namespace viewer.
type Namespace struct {
	ResourceViewer
	conn *grpc.ClientConn
}

// NewNamespace returns a new viewer.
func NewNamespace() ResourceViewer {
	var n Namespace
	conn, err := grpc.Dial(w2Server, grpc.WithInsecure())
	if err != nil {
		log.Panic().Msgf("connect grpc server failed:", err)
	}
	n.conn = conn
	fn := func() dao.MyFactory { return watch.NewProcessFactory(n.conn) }
	// factory := watch.NewProcessFactory(n.conn)
	n.ResourceViewer = NewBrowser("namespace", fn)

	// n.GetTable().SetDecorateFn(n.decorate)
	n.GetTable().SetColorerFn(render.Namespace{}.ColorerFunc())
	n.GetTable().SetEnterFn(n.switchNs)
	n.AddBindKeysFn(n.bindKeys)

	return &n
}

func (n *Namespace) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyU:      ui.NewKeyAction("Use", n.useNsCmd, true),
		ui.KeyShiftS: ui.NewKeyAction("Sort Status", n.GetTable().SortColCmd(statusCol, true), false),
	})
}

func (n *Namespace) switchNs(app *App, model ui.Tabular, gvr, path string) {
	ns := n.GetTable().GetSelectedItem2()
	n.useNamespace(ns)
	app.gotoResource("process", "", false)
}

func (n *Namespace) useNsCmd(evt *tcell.EventKey) *tcell.EventKey {
	path := n.GetTable().GetSelectedItem2()
	if path == "" {
		return nil
	}
	n.useNamespace(path)

	return nil
}

func (n *Namespace) useNamespace(fqn string) {
	_, ns := client.Namespaced(fqn)
	if err := n.App().switchNS(ns); err != nil {
		n.App().Flash().Err(err)
		return
	}
	if err := n.App().Config.SetActiveNamespace(ns); err != nil {
		n.App().Flash().Err(err)
		return
	}

	n.App().Flash().Infof("Namespace %s is now active!", ns)
	if err := n.App().Config.Save(); err != nil {
		log.Error().Err(err).Msg("Config file save failed!")
	}
}

func (n *Namespace) decorate(data render.TableData) render.TableData {
	if len(data.RowEvents) == 0 {
		return data
	}

	// checks if all ns is in the list if not add it.
	if _, ok := data.RowEvents.FindIndex(client.NamespaceAll); !ok {
		data.RowEvents = append(data.RowEvents,
			render.RowEvent{
				Kind: render.EventUnchanged,
				Row: render.Row{
					ID:     client.NamespaceAll,
					Fields: render.Fields{client.NamespaceAll, "Active", "", "", time.Now().String()},
				},
			},
		)
	}

	for _, re := range data.RowEvents {
		if config.InList(n.App().Config.FavNamespaces(), re.Row.ID) {
			re.Row.Fields[0] += favNSIndicator
			re.Kind = render.EventUnchanged
		}
		if n.App().Config.ActiveNamespace() == re.Row.ID {
			re.Row.Fields[0] += defaultNSIndicator
			re.Kind = render.EventUnchanged
		}
	}

	return data
}
