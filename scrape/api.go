// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ProcessHandler(w http.ResponseWriter, req *http.Request) {
	p := CollectorMap["process"]
	err := p.Start()
	if err != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("peek error happened: %s", err))
		return
	}

	var ns NameSpace
	ns.Erect(p.Peek())
	r := new(Report)
	output, err := r.Render(ns.Groups())
	if err != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("render error happened: %s", err))
	}
	_, _ = io.WriteString(w, output)
}

func GetProcesses(w http.ResponseWriter, req *http.Request) {
	p := CollectorMap["process"]
	err := p.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("peek error happened: %s", err))
		return
	}

	p.Refresh()

	var ns NameSpace
	ns.Erect(p.Peek())

	g := ns.Groups()
	content, err := json.Marshal(g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("json mashal error : %s", err))
		return
	}

	_, _ = w.Write(content)
}
