// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package scrape

import (
	"github.com/liupzmin/weewoe/internal/render"
)

// Collector 收集数据
type Collector interface {
	// Start 开启后台收集任务
	Start() error
	// Stop 停止收集
	Stop()
	// Refresh 负责手动刷新
	Refresh()
	// Peek Peek返回未加工的数据
	Peek() render.Rows
	Running() bool
	// AddListener 注册数据订阅者
	AddListener(TableListener)
	// RemoveListener 移除订阅者
	RemoveListener(TableListener)
	SendCommand(int64)
}

type TableListener interface {
	// TableDataChanged notifies the model data changed.
	TableDataChanged(rows render.Rows)

	// TableLoadFailed notifies the load failed.
	TableLoadFailed(error)

	Chan() chan []byte
}
