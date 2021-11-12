package mailer

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/bigmate/idm/pkg/logger"
	"github.com/bigmate/notification/internal/config"
	mail "github.com/xhit/go-simple-mail/v2"
)

//Mailer is general purpose mailer service
type Mailer interface {
	Send(ctx context.Context, options ...Option) error
}

//parameters is the email parameters struct
type parameters struct {
	sender         string
	receiver       string
	subject        string
	template       io.Reader
	connectTimeout time.Duration
	sendTimeout    time.Duration
}

//Option is an option function
type Option func(p *parameters)

//WithSender sets sender
func WithSender(sender string) Option {
	return func(p *parameters) {
		p.sender = sender
	}
}

//WithReceiver sets receiver
func WithReceiver(receiver string) Option {
	return func(p *parameters) {
		p.receiver = receiver
	}
}

//WithSubject sets subject of an email
func WithSubject(subject string) Option {
	return func(p *parameters) {
		p.subject = subject
	}
}

//WithConnectTimeout sets connect timeout
func WithConnectTimeout(timeout time.Duration) Option {
	return func(p *parameters) {
		p.connectTimeout = timeout
	}
}

//WithSendTimeout sets send timeout
func WithSendTimeout(timeout time.Duration) Option {
	return func(p *parameters) {
		p.sendTimeout = timeout
	}
}

//WithTemplate sets template
func WithTemplate(template io.Reader) Option {
	return func(p *parameters) {
		p.template = template
	}
}

//mailer is an implementation of Mailer interface
type mailer struct {
	sender                string
	host                  string
	port                  int
	username              string
	password              string
	encryption            mail.Encryption
	defaultConnectTimeout time.Duration
	defaultSendTimeout    time.Duration
}

func NewMailer(config *config.Config) Mailer {
	return &mailer{
		sender:                config.Smtp.Sender,
		host:                  config.Smtp.Host,
		port:                  config.Smtp.Port,
		username:              config.Smtp.Username,
		password:              config.Smtp.Password,
		defaultConnectTimeout: time.Minute,
		defaultSendTimeout:    time.Minute * 5,
		encryption:            mail.EncryptionSTARTTLS,
	}
}

func (m *mailer) Send(_ context.Context, options ...Option) error {
	p := &parameters{
		sender:         m.sender,
		connectTimeout: m.defaultConnectTimeout,
		sendTimeout:    m.defaultSendTimeout,
		template:       &bytes.Buffer{},
	}

	for _, option := range options {
		option(p)
	}

	if p.receiver == "" {
		panic("receiver is not set")
	}

	server := mail.NewSMTPClient()
	server.Host = m.host
	server.Port = m.port
	server.Username = m.username
	server.Password = m.password
	server.Encryption = m.encryption
	server.KeepAlive = false
	server.ConnectTimeout = p.connectTimeout
	server.SendTimeout = p.sendTimeout

	smtpClient, err := server.Connect()
	if err != nil {
		logger.Errorf("mailer: failed to connect: %v", err)
		return err
	}

	email := mail.NewMSG()
	email = email.SetFrom(p.sender)
	email = email.AddTo(p.receiver)
	email = email.SetSubject(p.subject)

	emailBytes, err := io.ReadAll(p.template)
	if err != nil {
		logger.Errorf("mailer: failed to read template: %v", err)
		return err
	}

	email.SetBodyData(mail.TextHTML, emailBytes)

	err = email.Send(smtpClient)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
