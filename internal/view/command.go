package view

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/model"
	"github.com/rs/zerolog/log"
)

var (
	customViewers MetaViewers
)

// Command represents a user command.
type Command struct {
	app *App

	mx sync.Mutex
}

// NewCommand returns a new command.
func NewCommand(app *App) *Command {
	return &Command{
		app: app,
	}
}

// Init initializes the command.
func (c *Command) Init() error {

	customViewers = loadCustomViewers()

	return nil
}

// Exec the Command by showing associated display.
func (c *Command) run(cmd, path string, clearStack bool) error {
	if c.specialCmd(cmd, path) {
		return nil
	}
	cmds := strings.Split(cmd, " ")
	gvr, v, err := c.viewMetaFor(cmds[0])
	if err != nil {
		return err
	}

	// checks if Command includes a namespace
	ns := c.app.Config.ActiveNamespace()
	if len(cmds) == 2 {
		ns = cmds[1]
	}
	if err := c.app.switchNS(ns); err != nil {
		return err
	}

	return c.exec(cmd, gvr, c.componentFor(gvr, path, v), clearStack)

}

func (c *Command) defaultCmd() error {
	var view string
	// view := c.app.Config.ActiveView()
	if view == "" {
		return c.run("process", "", true)
	}
	tokens := strings.Split(view, " ")
	cmd := view
	if len(tokens) == 1 {
		if !isContextCmd(tokens[0]) {
			cmd = tokens[0] + " " + c.app.Config.ActiveNamespace()
		}
	}

	if err := c.run(cmd, "", true); err != nil {
		log.Error().Err(err).Msgf("Default run command failed %q", cmd)
		c.app.cowCmd(err.Error())
		return err
	}
	return nil
}

func isContextCmd(c string) bool {
	return c == "ctx" || c == "context"
}

func (c *Command) specialCmd(cmd, path string) bool {
	cmds := strings.Split(cmd, " ")
	switch cmds[0] {
	case "cow":
		c.app.cowCmd(path)
		return true
	case "q", "Q", "quit":
		c.app.BailOut()
		return true
	case "?", "h", "help":
		c.app.helpCmd(nil)
		return true
	default:
	}
	return false
}

func (c *Command) viewMetaFor(cmd string) (string, *MetaViewer, error) {

	v, ok := customViewers[cmd]
	if !ok {
		return "", nil, fmt.Errorf("`%s` viewer not found", cmd)
	}

	return cmd, &v, nil
}

func (c *Command) componentFor(cat, path string, v *MetaViewer) ResourceViewer {
	var view ResourceViewer
	if v.viewerFn != nil {
		view = v.viewerFn()
	} else {
		// todo: 补充默认的 factory
		view = NewBrowser(cat, nil)
	}

	// 集群信息订阅 Table 数据
	if cat == "process" {
		view.GetTable().GetModel().AddListener(c.app.clusterInfo())
	}

	view.SetInstance(path)
	if v.enterFn != nil {
		view.GetTable().SetEnterFn(v.enterFn)
	}

	return view
}

func (c *Command) exec(cmd, gvr string, comp model.Component, clearStack bool) (err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		log.Error().Msgf("Something bad happened! %#v", e)
	//		c.app.Content.Dump()
	//		log.Debug().Msgf("History %v", c.app.cmdHistory.List())
	//
	//		hh := c.app.cmdHistory.List()
	//		if len(hh) == 0 {
	//			_ = c.run("process", "", true)
	//		} else {
	//			_ = c.run(hh[0], "", true)
	//		}
	//		err = fmt.Errorf("Invalid command %q", cmd)
	//	}
	//}()

	if comp == nil {
		return fmt.Errorf("No component found for %s", gvr)
	}
	c.app.Flash().Infof("Viewing %s...", client.NewGVR(gvr).R())
	if tokens := strings.Split(cmd, " "); len(tokens) >= 2 {
		cmd = tokens[0]
	}
	c.app.Config.SetActiveView(cmd)
	if err := c.app.Config.Save(); err != nil {
		log.Error().Err(err).Msg("Config save failed!")
	}
	if clearStack {
		c.app.Content.Stack.Clear()
	}

	if err := c.app.inject(comp); err != nil {
		return err
	}
	c.app.cmdHistory.Push(cmd)

	return
}
