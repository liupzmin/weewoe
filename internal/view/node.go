package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/weewoe/internal/ui"
)

// Node represents a node view.
type Node struct {
	ResourceViewer
}

// NewNode returns a new node view.
func NewNode(cat string) ResourceViewer {
	n := Node{
		// todo: 待实现
		ResourceViewer: NewBrowser(cat, nil),
	}
	n.AddBindKeysFn(n.bindKeys)
	n.GetTable().SetEnterFn(n.showProcess)

	return &n
}

func (n *Node) bindDangerousKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		/*ui.KeyC: ui.NewKeyAction("Cordon", n.toggleCordonCmd(true), true),
		ui.KeyU: ui.NewKeyAction("Uncordon", n.toggleCordonCmd(false), true),
		ui.KeyR: ui.NewKeyAction("Drain", n.drainCmd, true),*/
	})
	cl := n.App().Config.K9s.CurrentCluster
	if n.App().Config.K9s.Clusters[cl].FeatureGates.NodeShell {
		aa.Add(ui.KeyActions{
			//ui.KeyS: ui.NewKeyAction("Shell", n.sshCmd, true),
		})
	}
}

func (n *Node) bindKeys(aa ui.KeyActions) {
	aa.Delete(ui.KeySpace, tcell.KeyCtrlSpace)

	if !n.App().Config.K9s.IsReadOnly() {
		n.bindDangerousKeys(aa)
	}

	aa.Add(ui.KeyActions{
		/*ui.KeyY:      ui.NewKeyAction("YAML", n.yamlCmd, true),
		ui.KeyShiftC: ui.NewKeyAction("Sort CPU", n.GetTable().SortColCmd(cpuCol, false), false),
		ui.KeyShiftM: ui.NewKeyAction("Sort MEM", n.GetTable().SortColCmd(memCol, false), false),*/
	})
}

func (n *Node) showProcess(a *App, _ ui.Tabular, _, path string) {}
