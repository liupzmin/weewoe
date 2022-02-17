package model

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liupzmin/weewoe/internal"
	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/dao"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
)

const REFRESH int64 = 1

// MyTable represents a table model.
type MyTable struct {
	cat         string // category
	namespace   string
	data        *render.TableData
	listeners   []TableListener
	inUpdate    int32
	instance    string
	mx          sync.RWMutex
	labelFilter string
	r           <-chan render.Rows
	once        sync.Once
}

// NewMyTable returns a new table model.
func NewMyTable(cat string) *MyTable {
	return &MyTable{
		cat:  cat,
		data: render.NewTableData(),
	}
}

// SetLabelFilter sets the labels filter.
func (t *MyTable) SetLabelFilter(f string) {
	t.mx.Lock()
	t.labelFilter = f
	t.mx.Unlock()
}

// SetInstance sets a single entry table.
func (t *MyTable) SetInstance(path string) {
	t.instance = path
}

// AddListener adds a new model listener.
func (t *MyTable) AddListener(l TableListener) {
	t.listeners = append(t.listeners, l)
}

// RemoveListener delete a listener from the list.
func (t *MyTable) RemoveListener(l TableListener) {
	victim := -1
	for i, lis := range t.listeners {
		if lis == l {
			victim = i
			break
		}
	}

	if victim >= 0 {
		t.mx.Lock()
		defer t.mx.Unlock()
		t.listeners = append(t.listeners[:victim], t.listeners[victim+1:]...)
	}
}

// Watch initiates model updates.
func (t *MyTable) Watch(ctx context.Context) error {
	// t.reload(ctx)
	go t.updater(ctx)

	return nil
}

// Refresh updates the table content.
func (t *MyTable) Refresh(ctx context.Context) error {
	t.reload(ctx)
	return nil
}

// Get returns a resource instance if found, else an error.
func (t *MyTable) Get(ctx context.Context, path string) (runtime.Object, error) {
	return nil, nil
}

// Delete deletes a resource.
func (t *MyTable) Delete(ctx context.Context, path string, cascade, force bool) error {
	return nil
}

// GetNamespace returns the model namespace.
func (t *MyTable) GetNamespace() string {
	return t.namespace
}

// SetNamespace sets up model namespace.
func (t *MyTable) SetNamespace(ns string) {
	t.namespace = ns
	t.data.Clear()
}

// InNamespace checks if current namespace matches desired namespace.
func (t *MyTable) InNamespace(ns string) bool {
	return len(t.data.RowEvents) > 0 && t.namespace == ns
}

// SetRefreshRate sets model refresh duration.
func (t *MyTable) SetRefreshRate(d time.Duration) {

}

// ClusterWide checks if resource is scope for all namespaces.
func (t *MyTable) ClusterWide() bool {
	return client.IsClusterWide(t.namespace)
}

// Empty returns true if no model data.
func (t *MyTable) Empty() bool {
	return len(t.data.RowEvents) == 0
}

// Count returns the row count.
func (t *MyTable) Count() int {
	return len(t.data.RowEvents)
}

// Peek returns model data.
func (t *MyTable) Peek() render.TableData {
	t.mx.RLock()
	defer t.mx.RUnlock()

	return t.data.Clone()
}

func (t *MyTable) openGate(ctx context.Context) <-chan render.Rows {
	t.once.Do(func() {
		factory, ok := ctx.Value(internal.KeyFactory).(dao.MyFactory)
		if !ok {
			log.Panic().Msgf("expected Factory in context but got %T", ctx.Value(internal.KeyFactory))
		}

		t.r = factory.Stream(t.cat)
	})
	return t.r
}

func (t *MyTable) updater(ctx context.Context) {
	defer log.Debug().Msgf("TABLE-UPDATER canceled -- %q", t.cat)

	r := t.openGate(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case oo := <-r:
			err := t.refresh(ctx, oo)
			if err != nil {
				log.Error().Err(err).Msgf("Refresh failed")
				t.fireTableLoadFailed(err)
				return
			}
		}
	}
}

func (t *MyTable) reload(ctx context.Context) {
	factory, ok := ctx.Value(internal.KeyFactory).(dao.MyFactory)
	if !ok {
		log.Panic().Msgf(""+
			"expected Factory in context but got %T", ctx.Value(internal.KeyFactory))
	}
	err := factory.SendCommand(REFRESH)
	if err != nil {
		log.Error().Msgf("Table reload failed: %s", err.Error())
	}
}

func (t *MyTable) refresh(ctx context.Context, rows render.Rows) error {
	if !atomic.CompareAndSwapInt32(&t.inUpdate, 0, 1) {
		log.Debug().Msgf("Dropping update...")
		return nil
	}
	defer atomic.StoreInt32(&t.inUpdate, 0)

	if err := t.reconcile(ctx, rows); err != nil {
		return err
	}
	t.fireTableChanged(t.Peek())

	return nil
}

func (t *MyTable) reconcile(ctx context.Context, rows render.Rows) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	meta := resourceMeta(t.cat)
	if t.labelFilter != "" {
		ctx = context.WithValue(ctx, internal.KeyLabels, t.labelFilter)
	}

	/*var rows render.Rows

	rows = make(render.Rows, 0)

	if err := meta.Renderer.Render(oo, t.namespace, &rows); err != nil {
		return err
	}*/

	// if labelSelector in place might as well clear the model data.
	/*sel, ok := ctx.Value(internal.KeyLabels).(string)
	if ok && sel != "" {
		t.data.Clear()
	}*/
	t.data.Clear()
	log.Info().Msgf("Begin to update: %+v", rows)
	t.data.Update(rows)
	t.data.SetHeader(t.namespace, meta.Renderer.Header(t.namespace))

	if len(t.data.Header) == 0 {
		return fmt.Errorf("fail to list resource %s", t.cat)
	}

	return nil
}

func (t *MyTable) fireTableChanged(data render.TableData) {
	t.mx.RLock()
	defer t.mx.RUnlock()

	for _, l := range t.listeners {
		l.TableDataChanged(data)
	}
}

func (t *MyTable) fireTableLoadFailed(err error) {
	for _, l := range t.listeners {
		l.TableLoadFailed(err)
	}
}

func resourceMeta(cat string) ResourceMeta {
	// todo: 处理不存在的情况
	meta, _ := Registry[cat]

	return meta
}
