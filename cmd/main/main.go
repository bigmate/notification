package main

import (
	"context"
	"emailservice/internal/pkg"

	"emailservice/pkg/logger"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hashicorp/go-multierror"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	apps, err := initApps(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	if err = run(ctx, apps...); err != nil {
		logger.Fatal(err)
	}
}

func run(ctx context.Context, apps ...pkg.App) error {
	var multiErrs error
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, a := range apps {
		wg.Add(1)
		app := a
		go func() {
			defer wg.Done()
			if err := app.Run(ctx); err != nil {
				multiErrs = multierror.Append(multiErrs, err)
				cancel()

			}
		}()
	}
	wg.Wait()
	return multiErrs
}

func initApps(ctx context.Context) ([]pkg.App, error) {
	//	conf, err := config.NewConfig()
	// if err != nil {
	//	return nil, err
	//	}
	logger.Info("config successfully initialized")

	return []pkg.App{}, nil
}
