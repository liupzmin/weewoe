// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of weewoe

package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"

	"github.com/liupzmin/weewoe/config"
	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/mail"
	pb "github.com/liupzmin/weewoe/proto"
	"github.com/liupzmin/weewoe/scrape"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	appName      = "w2psd"
	shortAppDesc = "A daemon for checking distributed processes' status."
	longAppDesc  = "w2psd is a daemon to check the status of variety of processes."
)

var (
	w2Flags *config.Flags

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Long:  longAppDesc,
		Run:   run,
	}
)

func init() {
	initFlags()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Panicf("run failed %s", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	log.SetLevel(*w2Flags.LogLevel)
	scrape.Init()
	cronSendMail()
	sendAlert()

	go httpServer()
	s := startGRPCServer()

	select {
	case <-waitSignals():
		log.Info("SHUTTING DOWN......")
		s.Stop()
		scrape.Stop()
	}
}

func startGRPCServer() *grpc.Server {
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

	return s
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
		ignoreTime := viper.GetStringSlice("alert.ignore_time")
		log.Debugf("The ignore time string is %v", ignoreTime)

		alert := scrape.Alert{URL: viper.GetString("alert.url"), IgnoreTimeString: ignoreTime}
		err := alert.ConvertTime(ignoreTime)
		if err != nil {
			log.Panicf("Parse alert.ignore_time error: %s", err)
		}

		scrape.CollectorMap["process"].AddListener(&alert)
	}
}

func httpServer() {
	http.HandleFunc("/", scrape.ProcessHandler)
	http.HandleFunc("/list", scrape.GetProcesses)
	if err := http.ListenAndServe(":9528", nil); err != nil {
		log.Panicf("http server start failed: %s", err.Error())
	}
}

func initFlags() {
	w2Flags = config.NewFlags()

	rootCmd.PersistentFlags().StringVarP(
		w2Flags.LogLevel,
		"logLevel", "l",
		config.DefaultLogLevel,
		"Specify a log level (info, warn, debug, trace, error)",
	)
	rootCmd.Flags().StringVarP(
		w2Flags.LogFile,
		"logFile", "",
		config.DefaultLogFile,
		"Specify the log file",
	)

	rootCmd.Flags()
}
