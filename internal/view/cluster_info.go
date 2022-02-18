package view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/tview"
	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/model"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/rs/zerolog/log"
)

// ClusterInfoListener registers a listener for model changes.
type ClusterInfoListener interface {
	// ClusterInfoChanged notifies the cluster meta was changed.
	ClusterInfoChanged(c model.Cluster)
}

// ClusterInfo represents a cluster info view.
type ClusterInfo struct {
	*tview.Table
	listeners []ClusterInfoListener
	app       *App
	styles    *config.Styles
}

// NewClusterInfo returns a new cluster info view.
func NewClusterInfo(app *App) *ClusterInfo {
	return &ClusterInfo{
		Table:  tview.NewTable(),
		app:    app,
		styles: app.Styles,
	}
}

// Init initializes the view.
func (c *ClusterInfo) Init() {
	c.SetBorderPadding(0, 0, 1, 0)
	c.app.Styles.AddListener(c)
	c.layout()
	c.StylesChanged(c.app.Styles)
}

// StylesChanged notifies skin changed.
func (c *ClusterInfo) StylesChanged(s *config.Styles) {
	c.styles = s
	c.SetBackgroundColor(s.BgColor())
	c.updateStyle()
}

func (c *ClusterInfo) TableDataChanged(data render.TableData) {
	sCol := data.Header.IndexOf("STATUS", false)
	pCol := data.Header.IndexOf("PORTS-STATUS", false)
	nCol := data.Header.IndexOf("NODE", false)
	if sCol != 2 || pCol != 5 || nCol != 3 {
		return
	}

	cl := c.count(data)
	c.fireChanged(cl)
	c.app.QueueUpdateDraw(func() {
		c.Clear()
		c.layout()
		// todo: 修改 domain
		row := c.setCell(0, fmt.Sprintf("[white::b]%s", "HFP"))
		row = c.setCell(row, fmt.Sprintf("[aqua::b]%d", cl.Total))
		row = c.setCell(row, fmt.Sprintf("[green::b]%d", cl.Healthy))
		row = c.setCell(row, fmt.Sprintf("[red::b]%d", cl.Killed))
		row = c.setCell(row, fmt.Sprintf("[red::b]%d", cl.PortClosed))
		row = c.setCell(row, fmt.Sprintf("[pink::b]%d", cl.Node))
		c.updateStyle()
	})
}

func (c *ClusterInfo) TableLoadFailed(_ error) {}

func (c *ClusterInfo) hasMetrics() bool {
	mx := c.app.Conn().HasMetrics()
	if mx {
		auth, err := c.app.Conn().CanI("", "metrics.k8s.io/v1beta1/nodes", client.ListAccess)
		if err != nil {
			log.Warn().Err(err).Msgf("No nodes metrics access")
		}
		mx = auth
	}

	return mx
}

// AddListener adds a new model listener.
func (c *ClusterInfo) AddListener(l ClusterInfoListener) {
	c.listeners = append(c.listeners, l)
}

// RemoveListener delete a listener from the list.
func (c *ClusterInfo) RemoveListener(l ClusterInfoListener) {
	victim := -1
	for i, lis := range c.listeners {
		if lis == l {
			victim = i
			break
		}
	}

	if victim >= 0 {
		c.listeners = append(c.listeners[:victim], c.listeners[victim+1:]...)
	}
}

func (c *ClusterInfo) layout() {
	for row, section := range []string{"Domain", "Total", "Healthy", "Killed", "Closed Ports", "Nodes"} {
		c.SetCell(row, 0, c.sectionCell(section))
		c.SetCell(row, 1, c.infoCell(render.NAValue))
	}
}

func (c *ClusterInfo) sectionCell(t string) *tview.TableCell {
	cell := tview.NewTableCell(t + ":")
	cell.SetAlign(tview.AlignLeft)
	cell.SetBackgroundColor(c.app.Styles.BgColor())

	return cell
}

func (c *ClusterInfo) infoCell(t string) *tview.TableCell {
	cell := tview.NewTableCell(t)
	cell.SetExpansion(2)
	cell.SetTextColor(c.styles.K9s.Info.FgColor.Color())
	cell.SetBackgroundColor(c.app.Styles.BgColor())

	return cell
}

func (c *ClusterInfo) setCell(row int, s string) int {
	if s == "" {
		s = render.NAValue
	}
	c.GetCell(row, 1).SetText(s)
	return row + 1
}

func (c *ClusterInfo) count(data render.TableData) model.Cluster {
	var total, healthy, killed, portClosed, node int
	total = len(data.RowEvents)
	hosts := make(map[string]struct{})
	for _, v := range data.RowEvents {
		switch v.Row.Fields[2] {
		case render.Running:
			healthy++
		case render.EXCEPTION:
			killed++
		}

		ports := strings.Split(v.Row.Fields[5], "|")
		for _, s := range ports {
			if s == render.CLOSED {
				portClosed++
			}
		}

		hosts[v.Row.Fields[3]] = struct{}{}
	}
	node = len(hosts)

	return model.Cluster{
		Total:      total,
		Healthy:    healthy,
		Killed:     killed,
		PortClosed: portClosed,
		Node:       node,
	}
}

func (c *ClusterInfo) fireChanged(cl model.Cluster) {
	for _, l := range c.listeners {
		l.ClusterInfoChanged(cl)
	}
}

const defconFmt = "%s %s level!"

func (c *ClusterInfo) setDefCon(cpu, mem int) {
	var set bool
	l := c.app.Config.K9s.Thresholds.LevelFor("cpu", cpu)
	if l > config.SeverityLow {
		c.app.Status(flashLevel(l), fmt.Sprintf(defconFmt, flashMessage(l), "CPU"))
		set = true
	}
	l = c.app.Config.K9s.Thresholds.LevelFor("memory", mem)
	if l > config.SeverityLow {
		c.app.Status(flashLevel(l), fmt.Sprintf(defconFmt, flashMessage(l), "Memory"))
		set = true
	}
	if !set && !c.app.IsBenchmarking() {
		c.app.ClearStatus(true)
	}
}

func (c *ClusterInfo) updateStyle() {
	for row := 0; row < c.GetRowCount(); row++ {
		c.GetCell(row, 0).SetTextColor(c.styles.K9s.Info.FgColor.Color())
		c.GetCell(row, 0).SetBackgroundColor(c.styles.BgColor())
		var s tcell.Style
		s = s.Bold(true)
		s = s.Foreground(c.styles.K9s.Info.SectionColor.Color())
		s = s.Background(c.styles.BgColor())
		c.GetCell(row, 1).SetStyle(s)
	}
}

// ----------------------------------------------------------------------------
// Helpers...

func flashLevel(l config.SeverityLevel) model.FlashLevel {
	// nolint:exhaustive
	switch l {
	case config.SeverityHigh:
		return model.FlashErr
	case config.SeverityMedium:
		return model.FlashWarn
	default:
		return model.FlashInfo
	}
}

func flashMessage(l config.SeverityLevel) string {
	// nolint:exhaustive
	switch l {
	case config.SeverityHigh:
		return "Critical"
	case config.SeverityMedium:
		return "Warning"
	default:
		return "OK"
	}
}
