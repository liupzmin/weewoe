package view

import (
	"fmt"
	"strings"

	"github.com/liupzmin/weewoe/internal/dao"

	"github.com/liupzmin/weewoe/internal/render"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/liupzmin/weewoe/internal/ui/dialog"
	"github.com/liupzmin/weewoe/internal/watch"
	"github.com/liupzmin/weewoe/log"

	"github.com/gdamore/tcell/v2"
	"google.golang.org/grpc"
)

// Process represents a process viewer.
type Process struct {
	ResourceViewer
	conn *grpc.ClientConn
}

// NewProcess returns a new viewer.
func NewProcess() ResourceViewer {
	var p Process

	// todo: 创建 grpc 连接
	conn, err := grpc.Dial("192.168.18.237:9527", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect grpc server failed: %s", err)
	}
	p.conn = conn

	fn := func() dao.MyFactory { return watch.NewProcessFactory(conn) }

	// factory := watch.NewProcessFactory(conn)
	p.ResourceViewer = NewBrowser("process", fn)

	p.AddBindKeysFn(p.bindKeys)
	p.GetTable().SetEnterFn(p.showTips)
	p.GetTable().SetColorerFn(render.Process{}.ColorerFunc())
	p.GetTable().SetDecorateFn(p.waitToDo)

	return &p
}

func (p *Process) waitToDo(data render.TableData) render.TableData {
	col := data.IndexOfHeader("PORTS-STATUS")
	for _, re := range data.RowEvents {
		pss := strings.Split(re.Row.Fields[col], "|")
		var s string
		for _, po := range pss {
			if po == render.CLOSED {
				s = fmt.Sprintf("%s|%s", s, "[red::b]"+po)
				continue
			}
			s = fmt.Sprintf("%s|%s", s, po)
		}
		re.Row.Fields[col] = strings.Trim(s, "|")
	}

	return data
}

func (p *Process) bindDangerousKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		tcell.KeyCtrlK: ui.NewKeyAction("Kill", p.killCmd, true),
	})
}

func (p *Process) bindKeys(aa ui.KeyActions) {
	if !p.App().Config.K9s.IsReadOnly() {
		p.bindDangerousKeys(aa)
	}

	aa.Add(ui.KeyActions{
		ui.KeyShiftR: ui.NewKeyAction("Sort Name", p.GetTable().SortColCmd("NAME", true), true),
		ui.KeyShiftS: ui.NewKeyAction("Sort Status", p.GetTable().SortColCmd(statusCol, true), true),
		ui.KeyShiftI: ui.NewKeyAction("Sort Start Time", p.GetTable().SortColCmd("START-TIME", true), true),
		ui.KeyShiftO: ui.NewKeyAction("Sort Node", p.GetTable().SortColCmd("NODE", true), true),
	})
	// aa.Add(resourceSorters(p.GetTable()))
}

func (p *Process) showTips(app *App, data ui.Tabular, gvr, path string) {
	// todo: 待编写回车处理
	/*var co model.Component
	if err := app.inject(co); err != nil {
		app.Flash().Err(err)
	}*/
	dialog.ShowBoard(p.App().Styles.Dialog(), p.App().Content.Pages, "敬请期待！")

}

// Handlers...

func (p *Process) killCmd(evt *tcell.EventKey) *tcell.EventKey {

	// todo: 待实现
	p.Refresh()

	return evt
}

func (p *Process) Stop() {
	p.ResourceViewer.Stop()
	_ = p.conn.Close()
}
