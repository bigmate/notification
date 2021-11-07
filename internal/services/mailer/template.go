package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/bigmate/notification/internal/models"
	"github.com/bigmate/notification/pkg/logger"
)

//go:embed templates
var f embed.FS

type TemplateFactory interface {
	PasswordReset(data models.PasswordReset) (io.Reader, error)
}

type tmplFactory struct {
	tmpl *template.Template
}

func NewTemplateFactory() (TemplateFactory, error) {
	tmpl, err := template.ParseFS(f)
	if err != nil {
		return nil, err
	}
	return &tmplFactory{tmpl: tmpl}, nil
}

func (t *tmplFactory) PasswordReset(data models.PasswordReset) (io.Reader, error) {
	buf := &bytes.Buffer{}
	if err := t.tmpl.ExecuteTemplate(buf, "password-reset", data); err != nil {
		logger.Errorf("mailer: failed to construct")
		return nil, err
	}
	return buf, nil
}
