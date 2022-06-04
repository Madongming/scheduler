package store

type JobList struct {
	Concurrency       int      `json:"concurrency"`
	MaxScheduledCount int      `json:"maxScheduledCount"`
	MaxHistory        int      `json:"maxHistory"`
	Jobs              []string `json:"jobs"`
}

type Job struct {
	Name             string `json:"name"`
	Display          string `json:"display"`
	Type             string `json:"type"`
	ScheduleDuration int64  `json:"scheduleDuration"`
	Timeout          int64  `json:"timeout"`
	RetryTimes       uint8  `json:"retryTimes"`
	RetryWait        int64  `json:"retryWait"`
}
