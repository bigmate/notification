package background

import "context"

//Job is the job interface
type Job interface {
	Name() string
	Executable() func(ctx context.Context) error
}

//job is the implementation of the Job interface
type job struct {
	name string
	exe  func(ctx context.Context) error
}

func NewJob(name string, exe func(ctx context.Context) error) Job {
	return &job{
		name: name,
		exe:  exe,
	}
}

func (j *job) Name() string {
	return j.name
}

func (j *job) Executable() func(ctx context.Context) error {
	return j.exe
}
