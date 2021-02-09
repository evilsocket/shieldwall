package mailer

import (
	"github.com/evilsocket/islazy/log"
	"gopkg.in/gomail.v2"
	"sync"
)

type Mailer struct {
	sync.Mutex
	conf   Config
	dialer *gomail.Dialer
}

func New(conf Config) (*Mailer, error) {
	return &Mailer{
		conf:   conf,
		dialer: gomail.NewDialer(conf.Address, conf.Port, conf.Username, conf.Password),
	}, nil
}

func (m *Mailer) Send(from, to, subject, body string) error {
	m.Lock()
	defer m.Unlock()

	log.Debug("sending email to %s via %s:%d ...", to, m.conf.Address, m.conf.Port)

	msg := gomail.NewMessage()

	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}