package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/bigmate/idm/pkg/app"
	"github.com/bigmate/notification/internal/config"
	"github.com/bigmate/notification/internal/services/kafka"
	"github.com/bigmate/notification/internal/services/mailer"

	"github.com/bigmate/idm/pkg/logger"
	"github.com/bigmate/notification/internal/services/background"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	apps, err := initApps()
	if err != nil {
		logger.Fatalf("failed to init apps: %v", err)
	}

	runner := app.NewRunner(apps...)

	if err = runner.Run(ctx); err != nil {
		logger.Fatalf("failed to start app: %v", err)
	}
}

func initApps() ([]app.App, error) {
	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	logger.Info("config successfully initialized")

	bgScheduler := background.NewService()
	mailService := mailer.NewMailer(conf)
	kService, err := kafka.NewService(mailService, bgScheduler)

	if err != nil {
		return nil, err
	}

	logger.Info("kafka successfully initialized")

	return []app.App{bgScheduler, kService}, nil
}
