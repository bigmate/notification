package background

import (
	"context"

	"github.com/bigmate/notification/internal/pkg"
	"github.com/bigmate/notification/pkg/logger"
)

//Service is the background service interface
type Service interface {
	pkg.App
	Schedule(job Job)
}

//service is the service interface implementation
type service struct {
	ctx context.Context
}

func NewService() Service {
	return &service{
		ctx: context.Background(),
	}
}

//Schedule schedules a job
func (s *service) Schedule(job Job) {
	go func() {
		//TODO: add tracing
		logger.Infof("starting job %s", job.Name())
		defer logger.Infof("finished %s", job.Name())

		execute := job.Executable()
		if err := execute(s.ctx); err != nil {
			logger.Errorf("failed to execute job %s, %v", job.Name(), err)
		}
	}()
}

//Run runs a job
func (s *service) Run(ctx context.Context) error {
	s.ctx = ctx
	<-ctx.Done()
	return nil
}