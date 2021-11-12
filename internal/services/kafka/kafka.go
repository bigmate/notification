package kafka

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/bigmate/idm/pkg/app"
	"github.com/bigmate/idm/pkg/logger"
	"github.com/bigmate/notification/internal/models"
	"github.com/bigmate/notification/internal/services/background"
	"github.com/bigmate/notification/internal/services/mailer"
)

//service is a long-running kafka topics consumer
type service struct {
	mail     mailer.Mailer
	tmpl     mailer.TemplateFactory
	bg       background.Service
	consumer sarama.Consumer
}

func NewService(mail mailer.Mailer, bg background.Service) (app.App, error) {
	conf := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{""}, conf)
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
		bg:       bg,
	}, nil
}

//Run implements app.App interface
func (s *service) Run(ctx context.Context) error {
	if err := s.listenTopic("sign-up", s.processSignUp); err != nil {
		return err
	}

	if err := s.listenTopic("password-reset", s.processPasswordReset); err != nil {
		return err
	}

	<-ctx.Done()

	if err := s.consumer.Close(); err != nil {
		logger.Errorf("kafka: failed to close consumer: %v", err)
	}

	return nil
}

type processor func(ctx context.Context, message []byte) error

func (s *service) listenTopic(topic string, process processor) error {
	partitions, err := s.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, conErr := s.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if conErr != nil {
			return conErr
		}

		go func(consumer sarama.PartitionConsumer) {
			for message := range consumer.Messages() {
				job := background.NewJob(
					fmt.Sprintf("%s-consumer", message.Topic),
					func(ctx context.Context) error {
						return process(ctx, message.Value)
					},
				)
				s.bg.Schedule(job)
			}
		}(pc)
	}

	return nil
}

func (s *service) processSignUp(ctx context.Context, msg []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(msg))
	md := models.Signup{}

	if decodeErr := decoder.Decode(&md); decodeErr != nil {
		logger.Errorf("kafka: failed to decode %v", decodeErr)
		return decodeErr
	}

	r, tmplErr := s.tmpl.Signup(md)
	if tmplErr != nil {
		logger.Errorf("kafka: failed to create signup template %v", tmplErr)
		return tmplErr
	}

	return s.mail.Send(ctx,
		mailer.WithReceiver(md.Receiver),
		mailer.WithSubject("Welcome to Go Classifieds"),
		mailer.WithTemplate(r),
	)
}

func (s *service) processPasswordReset(ctx context.Context, msg []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(msg))
	md := models.PasswordReset{}

	if decodeErr := decoder.Decode(&md); decodeErr != nil {
		logger.Errorf("kafka: failed to decode %v", decodeErr)
		return decodeErr
	}

	r, tmplErr := s.tmpl.PasswordReset(md)
	if tmplErr != nil {
		logger.Errorf("kafka: failed to create password-reset template %v", tmplErr)
		return tmplErr
	}

	return s.mail.Send(ctx,
		mailer.WithReceiver(md.Receiver),
		mailer.WithSubject("Password reset"),
		mailer.WithTemplate(r),
	)
}
