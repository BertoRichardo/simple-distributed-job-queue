package mutation

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"
)

type JobMutation struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

func (q *JobMutation) Enqueue(ctx context.Context, args struct{ Task string }) (*resolver.JobResolver, error) {
	job, err := q.jobService.Enqueue(ctx, args.Task)
	if err != nil {
		return nil, err
	}
	return &resolver.JobResolver{Data: *job}, nil
}

func (q *JobMutation) SimultaneousCreateJob(ctx context.Context, args struct {
	Count      float64
	TaskPrefix string
}) (*[]*resolver.JobResolver, error) {
	jobs, err := q.jobService.SimultaneousCreateJob(ctx, int(args.Count), args.TaskPrefix)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*resolver.JobResolver, len(jobs))
	for i, job := range jobs {
		resolvers[i] = &resolver.JobResolver{Data: *job}
	}
	return &resolvers, nil
}

func (q *JobMutation) SimulateUnstableJob(ctx context.Context) (*resolver.JobResolver, error) {
	job, err := q.jobService.SimulateUnstableJob(ctx)
	if err != nil {
		return nil, err
	}
	return &resolver.JobResolver{Data: *job}, nil
}

func (q *JobMutation) CreateManyUnstableJobs(ctx context.Context, args struct{ Count float64 }) (*[]*resolver.JobResolver, error) {
	jobs, err := q.jobService.CreateManyUnstableJobs(ctx, int(args.Count))
	if err != nil {
		return nil, err
	}
	resolvers := make([]*resolver.JobResolver, len(jobs))
	for i, job := range jobs {
		resolvers[i] = &resolver.JobResolver{Data: *job}
	}
	return &resolvers, nil
}

func NewJobMutation(jobService _interface.JobService, dataloader *_dataloader.GeneralDataloader) JobMutation {
	return JobMutation{
		jobService: jobService,
		dataloader: dataloader,
	}
}
