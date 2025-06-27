// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package scrape

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/liupzmin/weewoe/util/xtime"

	"github.com/go-resty/resty/v2"
	"github.com/liupzmin/weewoe/internal/render"
	"github.com/liupzmin/weewoe/log"
)

type Message struct {
	Status      string      `json:"status"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}

type Labels struct {
	AlertName string `json:"alertname"`
	Service   string `json:"service"`
	Severity  string `json:"severity"`
	Instance  string `json:"instance"`
}

type Annotations struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type timePair struct {
	start time.Time
	end   time.Time
}

type Alert struct {
	URL              string
	IgnoreTimeString []string

	ignoreTime []timePair
	loc        *time.Location
}

func (a *Alert) TableDataChanged(rows render.Rows) {
	_ = a.ConvertTime(a.IgnoreTimeString)
	now := time.Now()
	log.Debugf("time now:%v", now)
	for _, tp := range a.ignoreTime {
		log.Debugf("start: %v, end: %v", tp.start, tp.end)
		if now.After(tp.start) && now.Before(tp.end) {
			log.Infof("shut up time！%v", a.IgnoreTimeString)
			return
		}
	}

	var ns NameSpace
	ns.Erect(rows)
	a.ProcessAlert(ns.Groups())
	a.PortAlert(ns.Groups())
}

func (a *Alert) TableLoadFailed(err error) {}

func (a *Alert) Chan() chan []byte { return nil }

func (a *Alert) ProcessAlert(data []Group) {
	msgs := make([]Message, 0)
	for _, v := range data {
		for _, p := range v.Processes {
			if p.State == Bad && !p.Suspend {
				msg := Message{
					Status: "firing",
					Labels: Labels{
						AlertName: p.Name,
						Service:   p.Name,
						Severity:  "critical",
						Instance:  p.Host,
					},
					Annotations: Annotations{
						Summary:     "Process Crash",
						Description: fmt.Sprintf("%s has crashed, please check!", v.Name),
					},
				}
				msgs = append(msgs, msg)
			}
		}
	}
	if len(msgs) > 0 {
		a.sent(msgs...)
	}
}

func (a *Alert) PortAlert(data []Group) {
	msgs := make([]Message, 0)
	for _, g := range data {
		for _, v := range g.Processes {
			for _, p := range v.Ports {
				if v.State == Good && p.State == 0 {
					msg := Message{
						Status: "firing",
						Labels: Labels{
							AlertName: v.Name,
							Service:   v.Name,
							Severity:  "critical",
							Instance:  v.Host,
						},
						Annotations: Annotations{
							Summary:     "port unreachable",
							Description: fmt.Sprintf("port %s cat't reach, please check!", p.Number),
						},
					}
					msgs = append(msgs, msg)
				}
			}
		}
	}
	if len(msgs) > 0 {
		a.sent(msgs...)
	}
}

type Result struct {
	Status string `json:"status"`
}

func (a *Alert) sent(msg ...Message) {
	result := new(Result)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		SetResult(result).
		Post(a.URL)

	if err != nil {
		log.Errorf("request for sending alert failed: %s", err.Error())
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Errorf("Sent alert, get wrong http code: %d", resp.StatusCode())
		return
	}

	if result.Status != "success" {
		log.Errorf("send alert failed: %s", result.Status)
		return
	}
}

func (a *Alert) ConvertTime(ts []string) error {
	if a.loc == nil {
		a.loc = xtime.GetUTC8Location()
	}
	var (
		tps []timePair
		now = time.Now()
		ymd = strings.Split(now.String(), " ")[0]
	)

	for _, tp := range ts {
		tpslice := strings.Split(tp, "-")
		if len(tpslice) != 2 {
			return fmt.Errorf("wrong alert ignore time config: %s", tp)
		}
		startTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", ymd, tpslice[0]), a.loc)
		if err != nil {
			return err
		}
		endTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", ymd, tpslice[1]), a.loc)
		if err != nil {
			return err
		}

		tps = append(tps, timePair{start: startTime, end: endTime})
	}
	a.ignoreTime = tps
	return nil
}
