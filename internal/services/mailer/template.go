package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/bigmate/idm/pkg/logger"
	"github.com/bigmate/notification/internal/models"
)

//go:embed templates
var f embed.FS

type TemplateFactory interface {
	PasswordReset(data models.PasswordReset) (io.Reader, error)
	Signup(data models.Signup) (io.Reader, error)
}

type tmplFactory struct {
	tmpl *template.Template
}

func NewTemplateFactory() (TemplateFactory, error) {
	tmpl, err := template.ParseFS(f, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}
	return &tmplFactory{tmpl: tmpl}, nil
}

//PasswordReset is the function that resets the password
func (t *tmplFactory) PasswordReset(data models.PasswordReset) (io.Reader, error) {
	buf := &bytes.Buffer{}
	if err := t.tmpl.ExecuteTemplate(buf, "password-reset.tmpl", data); err != nil {
		logger.Errorf("mailer: failed to construct %v", err)
		return nil, err
	}
	return buf, nil
}

//Signup is the function that sends the signup welcome email
func (t *tmplFactory) Signup(data models.Signup) (io.Reader, error) {
	buf := &bytes.Buffer{}
	if err := t.tmpl.ExecuteTemplate(buf, "signup.tmpl", data); err != nil {
		logger.Errorf("mailer: failed to construct template: %v", err)
		return nil, err
	}
	return buf, nil
}
