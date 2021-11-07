package mailer

import (
	"context"
	"io"
	"time"
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
	port                  string
	username              string
	password              string
	encryption            int
	keepAlive             bool
	defaultConnectTimeout time.Duration
	defaultSendTimeout    time.Duration
}

func NewMailer() Mailer {
	return nil
}

func (m *mailer) Send(ctx context.Context, options ...Option) error {
	p := &parameters{
		sender:         m.sender,
		connectTimeout: m.defaultConnectTimeout,
		sendTimeout:    m.defaultSendTimeout,
	}

	for _, option := range options {
		option(p)
	}

	if p.template == nil {
		panic("template is not set")
	}

	if p.receiver == "" {
		panic("receiver is not set")
	}

	return nil
}
