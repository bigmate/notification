package kafka

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/Shopify/sarama"
	"github.com/bigmate/notification/internal/models"
	"github.com/bigmate/notification/internal/pkg"
	"github.com/bigmate/notification/internal/services/mailer"
	"github.com/bigmate/notification/pkg/logger"
)

//service is the kafka service struct
type service struct {
	mail     mailer.Mailer
	tmpl     mailer.TemplateFactory
	consumer sarama.Consumer
}

//Run runs a job
func (s *service) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.consumer.Close()
	}()
	pc, err := s.consumer.ConsumePartition("signup", 1, 0)
	if err != nil {
		return err
	}

	messages := pc.Messages()

	for message := range messages {
		msgBytes := message.Value
		decoder := gob.NewDecoder(bytes.NewReader(msgBytes))
		md := models.Signup{}
		if err := decoder.Decode(&md); err != nil {
			logger.Errorf("kafka: failed to decode %v", err)
			continue
		}

		r, err := s.tmpl.Signup(md)
		if err != nil {
			logger.Errorf("kafka: failed to create signup template %v", err)
			continue
		}
		s.mail.Send(ctx,
			mailer.WithReceiver(md.Receiver),
			mailer.WithSubject("Welcome to Go Classifieds"),
			mailer.WithTemplate(r),
		)
	}
	return nil
}

func NewService(mail mailer.Mailer) (pkg.App, error) {
	consumer, err := sarama.NewConsumer([]string{""}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	tmpl, err := mailer.NewTemplateFactory()
	if err != nil {
		return nil, err
	}
	return &service{
		mail:     mail,
		consumer: consumer,
		tmpl:     tmpl,
	}, nil
}
