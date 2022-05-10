package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"

	"k8s.io/apimachinery/pkg/util/json"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/mail"
	pb "github.com/liupzmin/weewoe/proto"
	"github.com/liupzmin/weewoe/scrape"
	"github.com/liupzmin/weewoe/tmpl"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	scrape.Init()
	cronSendMail()
	sendAlert()

	go httpServer()

	s := grpc.NewServer()
	pb.RegisterStateServer(s, &scrape.State{})
	reflection.Register(s)

	go func() {
		lis, err := net.Listen("tcp", ":9527")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc serve failed: %v", err)
		}
	}()

	select {
	case <-waitSignals():
		log.Info("SHUTTING DOWN......")
		scrape.Stop()
		s.GracefulStop()
	}
}

func cronSendMail() {
	m := mail.New()

	go func() {
		for _ = range scrape.SendMail {
			m.Run()
		}
	}()

	if viper.GetBool("mail.send") {
		c := cron.New()
		_, err := c.AddJob(viper.GetString("mail.cron"), m)
		if err != nil {
			log.Panic("add job panic", log.FieldErr(err))
		}
		c.Start()
	}
}

func sendAlert() {
	if viper.GetBool("alert.notify") {
		alert := scrape.Alert{URL: viper.GetString("alert.url")}
		scrape.CollectorMap["process"].AddListener(alert)
	}
}

func httpServer() {
	http.HandleFunc("/", processHandler)
	http.HandleFunc("/list", getProcesses)
	if err := http.ListenAndServe(":9528", nil); err != nil {
		log.Panicf("http server start failed: %s", err.Error())
	}
}

func processHandler(w http.ResponseWriter, req *http.Request) {
	p := scrape.CollectorMap["process"]
	err := p.Start()
	if err != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("peek error happened: %s", err))
		return
	}

	var ns scrape.NameSpace
	ns.Erect(p.Peek())
	r := new(tmpl.Report)
	output, err := r.Render(ns.Groups())
	if err != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("render error happened: %s", err))
	}
	_, _ = io.WriteString(w, output)
}

func getProcesses(w http.ResponseWriter, req *http.Request) {
	p := scrape.CollectorMap["process"]
	err := p.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, fmt.Sprintf("peek error happened: %s", err))
		return
	}

	p.Refresh()

	var ns scrape.NameSpace
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
