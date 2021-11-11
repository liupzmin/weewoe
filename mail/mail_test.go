package mail

import (
	"os"
	"testing"
	"text/template"
)

func TestSendMail(t *testing.T) {
	data := Report{
		Title:              "HFP巡检报告",
		ReportName:         "HFP巡检报告",
		ReportDate:         "2021-11-11 07：00",
		Reporter:           "高运金融",
		NormalCount:        80,
		PlanStopCount:      0,
		PortExceptionCount: 0,
		PIDExceptionCount:  0,
		Groups: []Group{
			{
				Name: "Broker",
				Processes: []Process{
					{
						Name: "Authority",
						Host: "10.158.61.8",
						Ports: []Port{
							{
								Number: "20001",
								State:  1,
							},
							{
								Number: "20002",
								State:  0,
							},
						},
						State:     0,
						StartTime: "",
					},
					{
						Name: "Base",
						Host: "10.158.61.8",
						Ports: []Port{
							{
								Number: "20001",
								State:  1,
							},
						},
						State:     1,
						StartTime: "2021-11-11 05:30",
					},
				},
			},
			{
				Name: "Wechat",
				Processes: []Process{
					{
						Name:    "wechat",
						Host:    "10.158.61.8",
						Suspend: true,
						Ports: []Port{
							{
								Number: "8080",
								State:  0,
							},
						},
						State:     0,
						StartTime: "",
					},
					{
						Name: "wechat-web-admin",
						Host: "10.158.61.8",
						Ports: []Port{
							{
								Number: "20001",
								State:  1,
							},
						},
						State:     1,
						StartTime: "2021-11-11 05:30",
					},
					{
						Name: "wechat-web-client",
						Host: "10.158.61.8",
						Ports: []Port{
							{
								Number: "20001",
								State:  1,
							},
						},
						State:     1,
						StartTime: "2021-11-11 05:30",
					},
				},
			},
		},
	}

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
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

	tem, err := template.New("todos").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		t.Fatalf("temp parse error: %s", err.Error())
	}
	err = tem.Execute(os.Stdout, data)
	if err != nil {
		t.Fatalf("temp execute error: %s", err.Error())
	}
}
