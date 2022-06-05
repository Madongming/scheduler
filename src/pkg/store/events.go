package store

import "time"

type EventQueue struct {
	MaxStorage int           `json:"maxStorage"`
	MaxHistory int           `json:"maxHistory"`
	History    []JobInstance `json:"history"`
}

type JobInstance struct {
	Name      string    `json:"name"`
	Display   string    `json:"display"`
	TypeName  string    `json:"type"`
	State     string    `json:"state"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Done      bool      `json:"done"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
}

func NewJobInstance(job *Job) (JobInstance, error) {
	return JobInstance{
		Name:     job.Name,
		Display:  job.Display,
		TypeName: job.TypeName,
	}, nil
}
