package tmpl

import (
	"bytes"
	"text/template"
	"time"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/scrape"
	"github.com/spf13/viper"
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

func (r *Report) Render(gps []scrape.Group) (string, error) {
	r.suck(gps)
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

func (r *Report) suck(gps []scrape.Group) {
	r.Title = viper.GetString("report.title")
	r.ReportName = viper.GetString("report.title")
	r.ReportDate = time.Now().Format(scrape.TimeLayout)
	r.Reporter = viper.GetString("report.reporter")

	r.Groups = gps

	for _, v := range r.Groups {
		r.NormalCount += v.CountNormal()
		r.PlanStopCount += v.CountSuspend()
		r.PortExceptionCount += v.CountPortException()
		r.PIDExceptionCount += v.CountProcessException()
	}
}
