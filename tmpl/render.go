package tmpl

import (
	"bytes"
	"text/template"
	"time"

	"github.com/spf13/viper"

	"github.com/liupzmin/weewoe/scrape"

	"github.com/liupzmin/weewoe/log"
)

type Report struct {
	Title              string
	ReportName         string
	ReportDate         string
	Reporter           string
	NormalCount        int
	PlanStopCount      int
	PortExceptionCount int
	PIDExceptionCount  int
	Groups             []scrape.Group
}

func (r *Report) Render() (string, error) {
	r.suck()
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"isEven": func(i int) bool {
			if i%2 == 0 {
				return true
			}
			return false
		},
	}

	tem, err := template.New("todos").Funcs(funcMap).Parse(tm)
	if err != nil {
		log.DPanicf("parse template failed: %s", err.Error())
		return "", err
	}

	var buf bytes.Buffer
	err = tem.Execute(&buf, r)
	if err != nil {
		log.Errorf("template execute error :%s", err)
		return "", err
	}
	return buf.String(), nil
}

func (r *Report) suck() {
	r.Title = viper.GetString("report.title")
	r.ReportName = viper.GetString("report.title")
	r.ReportDate = time.Now().Format(scrape.TimeLayout)
	r.Reporter = viper.GetString("report.reporter")

	r.Groups = scrape.SC.MergeSort(scrape.SC.FetchPro(), scrape.SC.FetchPort())

	for _, v := range r.Groups {
		r.NormalCount += v.CountNormal()
		r.PlanStopCount += v.CountSuspend()
		r.PortExceptionCount += v.CountPortException()
		r.PIDExceptionCount += v.CountProcessException()
	}
}
