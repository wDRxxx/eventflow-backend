package smtp

import (
	"bytes"
	"html/template"
	"log/slog"
	"net"
	"net/smtp"
	"sync"

	"github.com/wDRxxx/eventflow-backend/internal/closer"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/mailer"
	"github.com/wDRxxx/eventflow-backend/internal/models"
)

type mail struct {
	wg             *sync.WaitGroup
	orderMailsChan chan *models.OrderMessage
	doneChan       chan struct{}

	login string
	host  string
	port  string

	auth      *smtp.Auth
	orderTmpl *template.Template
}

func NewSMTPMailer(config *config.MailerConfig, wg *sync.WaitGroup) (mailer.Mailer, error) {
	auth := smtp.PlainAuth("", config.Login(), config.Password(), config.Host())

	tmpl, err := template.ParseFiles(config.OrderTemplatePath())
	if err != nil {
		return nil, err
	}

	m := &mail{
		wg:             wg,
		orderMailsChan: make(chan *models.OrderMessage, 100),
		doneChan:       make(chan struct{}),
		login:          config.Login(),
		host:           config.Host(),
		port:           config.Port(),
		auth:           &auth,
		orderTmpl:      tmpl,
	}

	closer.Add(1, func() error {
		slog.Info("sending done signal to mail mailer...")
		m.doneChan <- struct{}{}

		return nil
	})

	closer.Add(2, func() error {
		slog.Info("closing api service channels...")
		close(m.orderMailsChan)
		close(m.doneChan)

		return nil
	})

	return m, nil
}

func (m *mail) Address() string {
	return net.JoinHostPort(m.host, m.port)
}

func (m *mail) ListenForMails() {
	for {
		select {
		case msg := <-m.orderMailsChan:
			m.wg.Add(1)
			go func() {
				defer m.wg.Done()

				err := m.sendOrderMail(msg)
				if err != nil {
					slog.Error("error sending order message: ", slog.Any("error", err), slog.Any("message", msg))
				}
			}()
		case <-m.doneChan:
			return
		}
	}
}

func (m *mail) SendOrderMail(msg *models.OrderMessage) {
	m.orderMailsChan <- msg
}

func (m *mail) sendOrderMail(msg *models.OrderMessage) error {
	var body bytes.Buffer

	err := m.orderTmpl.Execute(&body, msg)
	if err != nil {
		return err
	}

	data := []byte("From:  EventFlow\n" +
		"To: " + msg.To[0] + "\n" +
		"Subject: Your ticket to \"" + msg.EventTitle + "\"\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		body.String())

	err = m.SendHTMLMessage(data, msg.To)
	if err != nil {
		return err
	}

	return nil
}

func (m *mail) SendHTMLMessage(body []byte, to []string) error {
	err := smtp.SendMail(m.Address(), *m.auth, m.login, to, body)
	if err != nil {
		return err
	}

	return nil
}
