package scrape

import (
	"fmt"
	"net/http"

	"github.com/liupzmin/weewoe/log"

	"github.com/go-resty/resty/v2"

	"github.com/spf13/viper"
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

type Alert struct {
	URL string
}

func (a Alert) Notify() {
	if viper.GetBool("alert.notify") {
		go func() {
			for data := range SC.mch {
				a.ProcessAlert(data)
			}
		}()

		for data := range SC.nch {
			a.PortAlert(data)
		}
	}
}

func (a Alert) ProcessAlert(data []*ProcessState) {
	msgs := make([]Message, 0)
	for _, v := range data {
		if v.State == Bad && !v.Suspend {
			msg := Message{
				Status: "firing",
				Labels: Labels{
					AlertName: v.Name,
					Service:   v.Name,
					Severity:  "critical",
					Instance:  v.Host,
				},
				Annotations: Annotations{
					Summary:     "Process Crash",
					Description: fmt.Sprintf("%s has crashed, please check!", v.Name),
				},
			}
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) > 0 {
		a.sent(msgs...)
	}
}

func (a Alert) PortAlert(data []*PortState) {
	pros := SC.FetchPro()
	groups := SC.MergeSort(pros, data)
	msgs := make([]Message, 0)
	for _, g := range groups {
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

func (a Alert) sent(msg ...Message) {
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
