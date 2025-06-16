package entity

type JobStatusEnum string

const (
	StatusPending   JobStatusEnum = "PENDING"
	StatusRunning   JobStatusEnum = "RUNNING"
	StatusCompleted JobStatusEnum = "COMPLETED"
	StatusFailed    JobStatusEnum = "FAILED"
)

type Job struct {
	ID       string        `json:"id"`
	Task     string        `json:"task"`
	Status   JobStatusEnum `json:"status"` 
	Attempts int32         `json:"attempts"`
}

type JobStatus struct {
	Pending   int32 `json:"pending"`
	Running   int32 `json:"running"`
	Failed    int32 `json:"failed"`
	Completed int32 `json:"completed"`
}
