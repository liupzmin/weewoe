package main

import (
	"io"
	"net"
	"net/http"

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

	lis, err := net.Listen("tcp", ":9527")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	go func() {
		select {
		case <-waitSignals():
			s.Stop()
			scrape.GetCollector().Stop()
		}
	}()

	pb.RegisterStateServer(s, &scrape.State{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc serve failed: %v", err)
	}
}

func cronSendMail() {
	if viper.GetBool("mail.send") {
		m := mail.New()

		go func() {
			<-scrape.SendMailCH
			m.Run()
		}()

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
		go alert.Notify()
	}
}

func httpServer() {
	http.HandleFunc("/", processHandler)
	if err := http.ListenAndServe(":9528", nil); err != nil {
		log.Panicf("http server start failed: %s", err.Error())
	}
}

func processHandler(w http.ResponseWriter, req *http.Request) {
	r := new(tmpl.Report)
	output, err := r.Render()
	if err != nil {
		_, _ = io.WriteString(w, "error happened")
	}
	_, _ = io.WriteString(w, output)
}
