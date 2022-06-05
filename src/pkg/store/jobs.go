package store

import (
	"sync"
	"time"
)

var (
	mu sync.Mutex
)

type JobList struct {
	Concurrency       int      `json:"concurrency"`
	MaxScheduledCount int      `json:"maxScheduledCount"`
	MaxHistory        int      `json:"maxHistory"`
	Jobs              []string `json:"jobs"`
}

type Job struct {
	Name             string        `json:"name"`
	Display          string        `json:"display"`
	Type             string        `json:"type"`
	ScheduleDuration time.Duration `json:"scheduleDuration"`
	Timeout          time.Duration `json:"timeout"`
	RetryTimes       uint8         `json:"retryTimes"`
	RetryWait        time.Duration `json:"retryWait"`
}
