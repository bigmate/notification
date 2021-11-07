package pkg

import (
	"context"
)

//App is the application interface
type App interface {
	Run(ctx context.Context) error
}
