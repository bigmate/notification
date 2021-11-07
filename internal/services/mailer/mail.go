package mailer

import (
	"bytes"
	"context"
	"emailservice/pkg/logger"
	"embed"
	"fmt"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

//Mailer is the send email interface
type Mailer interface {
	Send(ctx context.Context, options ...Option) error
}

//Parameters is the email parameters struct
type Parameters struct {
	From           string
	To             string
	Subject        string
	Renderable     Renderable
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
}

//Option is the send options function
type Option func(p *Parameters)

//WithFrom sets from option as a parameter
func WithFrom(from string) Option {
	return func(p *Parameters) {
		p.From = from
	}
}

//WithTo sets to option as a parameter
func WithTo(to string) Option {
	return func(p *Parameters) {
		p.To = to
	}
}

//WithSubject sets subject option as a parameter
func WithSubject(subject string) Option {
	return func(p *Parameters) {
		p.Subject = subject
	}
}

//WithConnectTimeout sets subject option as a parameter
func WithConnectTimeout(timeout time.Duration) Option {
	return func(p *Parameters) {
		p.ConnectTimeout = timeout
	}
}

//WithConnectTimeout sets subject option as a parameter
func WithSendTimeout(timeout time.Duration) Option {
	return func(p *Parameters) {
		p.SendTimeout = timeout
	}
}

////WithRenderable sets from Renderable as a parameter
func WithRenderable(renderable Renderable) Option {
	return func(p *Parameters) {
		p.Renderable = renderable
	}
}

//mailer is the setup the mailer implementation
type mailer struct {
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

//go:embed templates
var emailTemplateFS embed.FS

func (m *mailer) Send(ctx context.Context, options ...Option) error {
	p := &Parameters{
		From:           "goclassified@mail.com",
		ConnectTimeout: m.defaultConnectTimeout,
		SendTimeout:    m.defaultSendTimeout,
	}

	for _, option := range options {
		option(p)
	}

	if p.Renderable == nil {
		panic("Renderable not set")
	}

	if p.To == "" {
		panic("No receiver email set")
	}
	return nil
}

func SendMail(ctx context.Context, from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)
	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		logger.Errorf("email: could not create template %s", err)
		return err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		//errorLog.Println(err)
		return err
	}

	formattedMsg := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)
	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)
	if err != nil {
		//app.errorLog.Println(err)
		return err
	}

	if err = t.ExecuteTemplate(&tpl, "body", data); err != nil {
		//app.errorLog.Println(err)
		return err
	}

	plainMsg := tpl.String()

	//send the email
	server := mail.NewSMTPClient()
	// server.Host = config.Config.Smtp.Host
	// server.Port = config.smtp.port
	// server.Username = config.smtp.username
	// server.Password = config.smtp.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from).AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMsg)
	email.AddAlternative(mail.TextPlain, plainMsg)

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	logger.Info("sent mail")
	return nil
}
