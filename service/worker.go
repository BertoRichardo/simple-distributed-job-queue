package service

import (
	"context"
	"jobqueue/entity"
	_interface "jobqueue/interface"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	unstableJobTaskName = "unstable-job"
	maxRetries          = 3
	retryDelay          = 2 * time.Second
)

type WorkerPool struct {
	repo     _interface.JobRepository
	jobQueue chan *entity.Job
	logger   *logrus.Logger
}

func NewWorkerPool(repo _interface.JobRepository, jobQueue chan *entity.Job, logger *logrus.Logger) *WorkerPool {
	return &WorkerPool{
		repo:     repo,
		jobQueue: jobQueue,
		logger:   logger,
	}
}

func (wp *WorkerPool) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go wp.worker(i + 1)
	}
	wp.logger.Infof("Worker Pool started with %d workers", workerCount)
}

func (wp *WorkerPool) worker(id int) {
	wp.logger.Infof("Worker %d started", id)
	for job := range wp.jobQueue {
		log := wp.logger.WithFields(logrus.Fields{"job_id": job.ID, "task": job.Task, "worker_id": id})
		log.Info("Picked up job")

		job.Status = entity.StatusRunning
		job.Attempts++
		wp.repo.Save(context.Background(), job)

		time.Sleep(1 * time.Second)

		isUnstableAndShouldFail := job.Task == unstableJobTaskName && job.Attempts < 3
		if isUnstableAndShouldFail {
			log.Warnf("Simulating failure, attempt %d", job.Attempts)

			if job.Attempts < maxRetries {
				log.Infof("Scheduling job for retry after %v", retryDelay)
				// re-queue to prefent deadlock
				go func(j *entity.Job) {
					time.Sleep(retryDelay)
					j.Status = entity.StatusPending

					// save PENDING and queue
					wp.repo.Save(context.Background(), j)
					wp.jobQueue <- j
				}(job)
			} else {
				log.Errorf("Job failed after max retries.")
				job.Status = entity.StatusFailed
				wp.repo.Save(context.Background(), job)
			}
		} else {
			log.Info("Job completed successfully")
			job.Status = entity.StatusCompleted
			wp.repo.Save(context.Background(), job)
		}
	}
}
