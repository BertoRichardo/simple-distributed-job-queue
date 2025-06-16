package _dataloader

import (
	"context"
	"jobqueue/entity"

	"github.com/graph-gophers/dataloader/v6"
)

func (s GeneralDataloader) JobBatchFunc(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	jobIDs := keys.Keys()
	results := make([]*dataloader.Result, len(jobIDs))

	jobs, err := s.jobRepo.FindManyByIDs(ctx, jobIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result{Error: err}
		}
		return results
	}

	jobMap := make(map[string]*entity.Job, len(jobs))
	for _, job := range jobs {
		jobMap[job.ID] = job
	}

	for i, key := range jobIDs {
		if job, ok := jobMap[key]; ok {
			results[i] = &dataloader.Result{Data: job}
		} else {
			results[i] = &dataloader.Result{Data: nil}
		}
	}

	return results
}