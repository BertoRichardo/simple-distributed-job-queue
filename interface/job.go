package _interface

import (
	"context"
	"jobqueue/entity"
)

type JobService interface {
	Enqueue(ctx context.Context, taskName string) (*entity.Job, error)
	GetAllJobs(ctx context.Context) ([]*entity.Job, error)
	SimultaneousCreateJob(ctx context.Context, count int, taskPrefix string) ([]*entity.Job, error)
	CreateManyUnstableJobs(ctx context.Context, count int) ([]*entity.Job, error)
	SimulateUnstableJob(ctx context.Context) (*entity.Job, error)
	GetJobByID(ctx context.Context, id string) (*entity.Job, error)
	GetJobStatus(ctx context.Context) (*entity.JobStatus, error) 
}

type JobRepository interface {
	Save(ctx context.Context, job *entity.Job) error
	FindByID(ctx context.Context, id string) (*entity.Job, error)
	FindAll(ctx context.Context) ([]*entity.Job, error)
	FindManyByIDs(ctx context.Context, ids []string) ([]*entity.Job, error)
}
