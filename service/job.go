package service

import (
	"context"
	"fmt"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"sync"

	"github.com/google/uuid"
)

type jobService struct {
	jobRepo  _interface.JobRepository
	jobQueue chan<- *entity.Job
}

func (s *jobService) Enqueue(ctx context.Context, taskName string) (*entity.Job, error) {
	newJob := &entity.Job{
		ID:       uuid.New().String(),
		Task:     taskName,
		Status:   entity.StatusPending,
		Attempts: 0,
	}
	if err := s.jobRepo.Save(ctx, newJob); err != nil {
		return nil, err
	}
	s.jobQueue <- newJob
	return newJob, nil
}

func (s *jobService) GetAllJobs(ctx context.Context) ([]*entity.Job, error) {
	return s.jobRepo.FindAll(ctx)
}

func (s *jobService) GetJobByID(ctx context.Context, id string) (*entity.Job, error) {
	return s.jobRepo.FindByID(ctx, id)
}

func (s *jobService) SimulateUnstableJob(ctx context.Context) (*entity.Job, error) {
	return s.Enqueue(ctx, unstableJobTaskName)
}

func (s *jobService) SimultaneousCreateJob(ctx context.Context, count int, taskPrefix string) ([]*entity.Job, error) {
	var wg sync.WaitGroup
	jobsChan := make(chan *entity.Job, count)
	errChan := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			job, err := s.Enqueue(ctx, fmt.Sprintf("%s-%d", taskPrefix, i))
			if err != nil {
				errChan <- err
				return
			}
			jobsChan <- job
		}(i)
	}

	wg.Wait()
	close(jobsChan)
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}
	createdJobs := make([]*entity.Job, 0, count)
	for job := range jobsChan {
		createdJobs = append(createdJobs, job)
	}
	return createdJobs, nil
}

func (s *jobService) GetJobStatus(ctx context.Context) (*entity.JobStatus, error) {
	allJobs, err := s.jobRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	stats := &entity.JobStatus{}
	for _, job := range allJobs {
		switch job.Status {
		case entity.StatusPending:
			stats.Pending++
		case entity.StatusRunning:
			stats.Running++
		case entity.StatusCompleted:
			stats.Completed++
		case entity.StatusFailed:
			stats.Failed++
		}
	}
	return stats, nil
}

func (s *jobService) CreateManyUnstableJobs(ctx context.Context, count int) ([]*entity.Job, error) {
	var wg sync.WaitGroup
	jobsChan := make(chan *entity.Job, count)
	errChan := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			job, err := s.Enqueue(ctx, "unstable-job")
			if err != nil {
				errChan <- err
				return
			}
			jobsChan <- job
		}()
	}

	wg.Wait()
	close(jobsChan)
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	createdJobs := make([]*entity.Job, 0, count)
	for job := range jobsChan {
		createdJobs = append(createdJobs, job)
	}
	return createdJobs, nil
}

type Initiator func(s *jobService) *jobService

func NewJobService() Initiator {
	return func(s *jobService) *jobService {
		return s
	}
}

func (i Initiator) SetJobRepository(jobRepository _interface.JobRepository) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobRepo = jobRepository
		return s
	}
}

func (i Initiator) SetJobQueue(jobQueue chan<- *entity.Job) Initiator {
	return func(s *jobService) *jobService {
		i(s).jobQueue = jobQueue
		return s
	}
}

func (i Initiator) Build() _interface.JobService {
	return i(&jobService{})
}
