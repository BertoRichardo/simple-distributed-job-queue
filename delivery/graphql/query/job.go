package query

import (
	"context"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/resolver"
	_interface "jobqueue/interface"
)

type JobQuery struct {
	jobService _interface.JobService
	dataloader *_dataloader.GeneralDataloader
}

func (q *JobQuery) Jobs(ctx context.Context) ([]*resolver.JobResolver, error) {
	jobs, err := q.jobService.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*resolver.JobResolver, len(jobs))
	for i, job := range jobs {
		resolvers[i] = &resolver.JobResolver{Data: *job}
	}
	return resolvers, nil
}

func (q *JobQuery) Job(ctx context.Context, args struct{ ID string }) (*resolver.JobResolver, error) {
	job, err := q.jobService.GetJobByID(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}
	return &resolver.JobResolver{Data: *job}, nil
}

func (q *JobQuery) JobStatus(ctx context.Context) (*resolver.JobStatusResolver, error) {
	stats, err := q.jobService.GetJobStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &resolver.JobStatusResolver{Data: *stats}, nil
}

func NewJobQuery(jobService _interface.JobService,
	dataloader *_dataloader.GeneralDataloader) JobQuery {
	return JobQuery{
		jobService: jobService,
		dataloader: dataloader,
	}
}
