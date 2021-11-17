package mail

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"sync"

	"github.com/liupzmin/weewoe/tmpl"

	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"

	"github.com/liupzmin/weewoe/log"
	"github.com/liupzmin/weewoe/scrape"
)

type Mail struct {
	From     string
	Passwd   string
	To       []string
	SMTPHost string
	SMTPPort int
}

func New() Mail {
	m := Mail{
		From:     viper.GetString("mail.user"),
		Passwd:   viper.GetString("mail.passwd"),
		To:       viper.GetStringSlice("mail.to"),
		SMTPHost: viper.GetString("mail.server"),
		SMTPPort: viper.GetInt("mail.port"),
	}
	return m
}

func (m Mail) Run() {
	m.knock()
	r := new(tmpl.Report)
	output, err := r.Render()
	if err != nil {
		return
	}
	m.Send(r.Title, output)
}

func (m Mail) Send(title, content string) {
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
}

// knock 通知 collector 收集进程数据并等待其完成后退出
func (m Mail) knock() {
	coll := scrape.GetCollector()
	pch := make(chan []*scrape.ProcessState, 1)
	porch := make(chan []*scrape.PortState, 1)

	var (
		i, j int
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		for p := range pch {
			_ = p
			i++
			if i == 1 {
				go coll.UnRegisterProChan(pch)
				wg.Done()
			}
		}
		log.Debugf("pch closed, bye bye")
	}()

	wg.Add(1)
	go func() {
		for p := range porch {
			_ = p
			j++
			if j == 1 {
				go coll.UnRegisterPortChan(porch)
				wg.Done()
			}
		}
		log.Debugf("porch closed, bye bye")
	}()

	coll.RegisterProChan(pch)
	coll.RegisterPortChan(porch)
	coll.ReCollect()
	wg.Wait()
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
