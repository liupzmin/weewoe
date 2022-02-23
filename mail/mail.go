package mail

import (
	"crypto/tls"
	"errors"
	"net/smtp"

	"github.com/liupzmin/weewoe/tmpl"

	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/scrape"
)

type Mail struct {
	scrape.Collector
	From     string
	Passwd   string
	To       []string
	SMTPHost string
	SMTPPort int
}

func New() *Mail {
	m := Mail{
		Collector: scrape.CollectorMap["process"],
		From:      viper.GetString("mail.user"),
		Passwd:    viper.GetString("mail.passwd"),
		To:        viper.GetStringSlice("mail.to"),
		SMTPHost:  viper.GetString("mail.server"),
		SMTPPort:  viper.GetInt("mail.port"),
	}
	return &m
}

func (m *Mail) Run() {
	rows := m.Peek()
	var ns scrape.NameSpace
	ns.Erect(rows)

	r := new(tmpl.Report)
	output, err := r.Render(ns.Groups())
	if err != nil {
		return
	}
	m.Send(r.Title, output)
}

func (m *Mail) Send(title, content string) {
	gm := gomail.NewMessage()

	gm.SetHeader("From", m.From)
	gm.SetHeader("To", m.To...)
	gm.SetHeader("Subject", title)
	gm.SetBody("text/html", content)

	auth := LoginAuth(m.From, m.Passwd)
	d := gomail.NewDialer(m.SMTPHost, m.SMTPPort, m.From, m.Passwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	d.Auth = auth

	if err := d.DialAndSend(gm); err != nil {
		log.Errorf("DialAndSend err %v:", err)
	}
	log.Info("Send Mail Success!")
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}
