package view

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/liupzmin/weewoe/internal"
	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/dao"
	"github.com/liupzmin/weewoe/internal/model"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/rs/zerolog/log"
)

// Browser represents a generic resource browser.
type Browser struct {
	*Table

	namespaces map[int]string
	factoryFn  dao.FactoryFn
	factory    dao.MyFactory
	contextFn  ContextFunc
	cancelFn   context.CancelFunc
	mx         sync.RWMutex
}

// NewBrowser returns a new browser.
func NewBrowser(cat string, fn dao.FactoryFn) ResourceViewer {
	return &Browser{
		Table:     NewTable(cat),
		factoryFn: fn,
	}
}

// Init watches all running pods in given namespace.
func (b *Browser) Init(ctx context.Context) error {

	if err := b.Table.Init(ctx); err != nil {
		return err
	}

	if b.App().IsRunning() {
		b.app.CmdBuff().Reset()
	}

	b.bindKeys(b.Actions())
	for _, f := range b.bindKeysFn {
		f(b.Actions())
	}

	// 设置统一的 namespace
	ns := client.CleanseNamespace(b.app.Config.ActiveNamespace())
	b.setNamespace(ns)
	row, _ := b.GetSelection()
	if row == 0 && b.GetRowCount() > 0 {
		b.Select(1, 0)
	}

	b.CmdBuff().SetSuggestionFn(b.suggestFilter())
	b.GetTable().Table.AddListener(b)

	return nil
}

// InCmdMode checks if prompt is active.
func (b *Browser) InCmdMode() bool {
	return b.CmdBuff().InCmdMode()
}

func (b *Browser) suggestFilter() model.SuggestionFunc {
	return func(s string) (entries sort.StringSlice) {
		if s == "" {
			if b.App().filterHistory.Empty() {
				return
			}
			return b.App().filterHistory.List()
		}

		s = strings.ToLower(s)
		for _, h := range b.App().filterHistory.List() {
			if s == h {
				continue
			}
			if strings.HasPrefix(h, s) {
				entries = append(entries, strings.Replace(h, s, "", 1))
			}
		}
		return
	}
}

func (b *Browser) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		tcell.KeyEscape: ui.NewSharedKeyAction("Filter Reset", b.resetCmd, false),
		tcell.KeyEnter:  ui.NewSharedKeyAction("Filter", b.filterCmd, false),
		tcell.KeyHelp:   ui.NewSharedKeyAction("Help", b.helpCmd, false),
	})
}

// SetInstance sets a single instance view.
func (b *Browser) SetInstance(path string) {
	b.GetModel().SetInstance(path)
}

// Start initializes browser updates.
func (b *Browser) Start() {
	b.factory = b.factoryFn()
	ns := b.app.Config.ActiveNamespace()
	if n := b.GetModel().GetNamespace(); !client.IsClusterScoped(n) {
		ns = n
	}
	if err := b.app.switchNS(ns); err != nil {
		log.Error().Err(err).Msgf("ns switch failed")
	}

	b.stop()
	b.GetModel().AddListener(b)
	b.Table.Start()
	b.CmdBuff().AddListener(b)
	if err := b.GetModel().Watch(b.prepareContext()); err != nil {
		b.App().Flash().Err(fmt.Errorf("Watcher failed for %s -- %w", b.Cat(), err))
	}
}

// Stop terminates browser updates.
func (b *Browser) Stop() {
	b.stop()
	b.Table.Table.Stop()
	b.factory.Terminate()
}

func (b *Browser) stop() {
	b.mx.Lock()
	{
		if b.cancelFn != nil {
			b.cancelFn()
			b.cancelFn = nil
		}
	}
	b.mx.Unlock()
	b.GetModel().RemoveListener(b)
	b.CmdBuff().RemoveListener(b)
	b.Table.Stop()
}

// BufferChanged indicates the buffer was changed.
func (b *Browser) BufferChanged(_, _ string) {}

// BufferCompleted indicates input was accepted.
func (b *Browser) BufferCompleted(text, _ string) {
	if ui.IsLabelSelector(text) {
		b.GetModel().SetLabelFilter(ui.TrimLabelSelector(text))
	} else {
		b.GetModel().SetLabelFilter("")
	}
}

// BufferActive indicates the buff activity changed.
func (b *Browser) BufferActive(state bool, k model.BufferKind) {
	if state {
		return
	}
	if err := b.GetModel().Refresh(b.prepareContext()); err != nil {
		log.Error().Err(err).Msgf("Refresh failed for %s", b.Cat())
	}
	b.app.QueueUpdateDraw(func() {
		b.Update(b.GetModel().Peek(), b.App().Conn().HasMetrics())
		if b.GetRowCount() > 1 {
			b.App().filterHistory.Push(b.CmdBuff().GetText())
		}
	})
}

func (b *Browser) prepareContext() context.Context {
	ctx := b.defaultContext()
	ctx, b.cancelFn = context.WithCancel(ctx)
	if b.contextFn != nil {
		ctx = b.contextFn(ctx)
	}
	if path, ok := ctx.Value(internal.KeyPath).(string); ok && path != "" {
		b.Path = path
	}

	return ctx
}

func (b *Browser) Refresh() {
	_ = b.GetModel().Refresh(b.prepareContext())
}

// Name returns the component name.
func (b *Browser) Name() string { return b.Cat() }

// SetContextFn populates a custom context.
func (b *Browser) SetContextFn(f ContextFunc) { b.contextFn = f }

// GetTable returns the underlying table.
func (b *Browser) GetTable() *Table { return b.Table }

// Aliases returns all available aliases.
func (b *Browser) Aliases() []string {
	// todo: 待处理
	return nil
}

// ----------------------------------------------------------------------------
// Model Protocol...

// TableDataChanged notifies view new data is available.
func (b *Browser) TableDataChanged(data render.TableData) {
	var cancel context.CancelFunc
	b.mx.RLock()
	cancel = b.cancelFn
	b.mx.RUnlock()

	if !b.app.ConOK() || cancel == nil || !b.app.IsRunning() {
		return
	}

	b.app.QueueUpdateDraw(func() {
		b.refreshActions()
		// todo: what is hasMetrics ?
		b.Update(data, false)
	})
}

// TableLoadFailed notifies view something went south.
func (b *Browser) TableLoadFailed(err error) {
	b.app.QueueUpdateDraw(func() {
		b.app.Flash().Err(err)
		b.App().ClearStatus(false)
	})
}

func (b *Browser) TableTick() {
	b.app.QueueUpdateDraw(func() {})
}

// ----------------------------------------------------------------------------
// Actions...

func (b *Browser) helpCmd(evt *tcell.EventKey) *tcell.EventKey {
	if b.CmdBuff().InCmdMode() {
		return nil
	}

	return evt
}

func (b *Browser) resetCmd(evt *tcell.EventKey) *tcell.EventKey {
	if !b.CmdBuff().InCmdMode() {
		b.CmdBuff().ClearText(false)
		return b.App().PrevCmd(evt)
	}

	b.CmdBuff().Reset()
	if ui.IsLabelSelector(b.CmdBuff().GetText()) {
		b.Start()
	}
	b.Refresh()

	return nil
}

func (b *Browser) filterCmd(evt *tcell.EventKey) *tcell.EventKey {
	if !b.CmdBuff().IsActive() {
		return evt
	}

	b.CmdBuff().SetActive(false)
	if ui.IsLabelSelector(b.CmdBuff().GetText()) {
		b.Start()
		return nil
	}
	b.Refresh()

	return nil
}

func (b *Browser) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	path := b.GetSelectedItem()
	if b.filterCmd(evt) == nil || path == "" {
		return nil
	}

	// todo： 处理没有设置 EnterFunc 的情况
	var f EnterFunc
	// f := describeResource
	if b.enterFn != nil {
		f = b.enterFn
	}
	f(b.app, b.GetModel(), b.Cat(), path)

	return nil
}

func (b *Browser) refreshCmd(*tcell.EventKey) *tcell.EventKey {
	b.app.Flash().Info("Refreshing...")
	b.Refresh()

	return nil
}

// switchNamespaceCmd 通过快捷键进入 namespace
func (b *Browser) switchNamespaceCmd(evt *tcell.EventKey) *tcell.EventKey {
	i, err := strconv.Atoi(string(evt.Rune()))
	if err != nil {
		log.Error().Err(err).Msgf("Fail to switch namespace")
		return nil
	}
	ns := strings.TrimSpace(b.namespaces[i])

	if err := b.app.switchNS(ns); err != nil {
		b.App().Flash().Err(err)
		return nil
	}
	b.setNamespace(ns)
	b.app.Flash().Infof("Viewing namespace `%s`...", ns)
	// 无需重新布局，仅刷新数据，通过设置 namespace 改变数据
	b.Refresh()
	b.UpdateTitle()
	b.SelectRow(1, true)
	b.app.CmdBuff().Reset()
	// 重新设置激活的 namespace 并保存配置
	if err := b.app.Config.SetActiveNamespace(b.GetModel().GetNamespace()); err != nil {
		log.Error().Err(err).Msg("Config save NS failed!")
	}
	if err := b.app.Config.Save(); err != nil {
		log.Error().Err(err).Msg("Config save failed!")
	}

	return nil
}

// ----------------------------------------------------------------------------
// Helpers...

func (b *Browser) setNamespace(ns string) {
	ns = client.CleanseNamespace(ns)
	if b.GetModel().InNamespace(ns) {
		return
	}

	b.GetModel().SetNamespace(ns)
}

func (b *Browser) defaultContext() context.Context {
	ctx := context.WithValue(context.Background(), internal.KeyFactory, b.factory)
	ctx = context.WithValue(ctx, internal.KeyGVR, b.Cat())
	if b.Path != "" {
		ctx = context.WithValue(ctx, internal.KeyPath, b.Path)
	}
	if ui.IsLabelSelector(b.CmdBuff().GetText()) {
		ctx = context.WithValue(ctx, internal.KeyLabels, ui.TrimLabelSelector(b.CmdBuff().GetText()))
	}
	ctx = context.WithValue(ctx, internal.KeyNamespace, client.CleanseNamespace(b.App().Config.ActiveNamespace()))

	return ctx
}

func (b *Browser) refreshActions() {
	if b.App().Content.Top().Name() != b.Name() {
		return
	}
	aa := ui.KeyActions{
		ui.KeyC:        ui.NewKeyAction("Copy", b.cpCmd, false),
		tcell.KeyEnter: ui.NewKeyAction("View", b.enterCmd, false),
		tcell.KeyCtrlR: ui.NewKeyAction("Refresh", b.refreshCmd, true),
	}

	b.namespaceActions(aa)

	// pluginActions(b, aa)
	// hotKeyActions(b, aa)
	for _, f := range b.bindKeysFn {
		f(aa)
	}
	b.Actions().Add(aa)
	b.app.Menu().HydrateMenu(b.Hints())
}

func (b *Browser) namespaceActions(aa ui.KeyActions) {
	if b.GetTable().Path != "" {
		return
	}
	b.namespaces = make(map[int]string, config.MaxFavoritesNS)
	aa[ui.Key0] = ui.NewKeyAction(client.NamespaceAll, b.switchNamespaceCmd, true)
	b.namespaces[0] = client.NamespaceAll
	index := 1
	// 从配置文件中拿到最爱的 namespace 并绑定键盘事件
	for _, ns := range b.app.Config.FavNamespaces() {
		if ns == client.NamespaceAll {
			continue
		}
		aa[ui.NumKeys[index]] = ui.NewKeyAction(ns, b.switchNamespaceCmd, true)
		b.namespaces[index] = ns
		index++
	}
}
