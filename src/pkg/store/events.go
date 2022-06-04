package store

type EventQueue struct {
	MaxStorage int           `json:"maxStorage"`
	MaxHistory int           `json:"maxHistory"`
	History    []JobInstance `json:"history"`
}

type JobInstance struct {
	Name           string `json:"name"`
	Display        string `json:"display"`
	Type           string `json:"type"`
	State          string `json:"state"`
	StartTimestamp int64  `json:"startTime"`
	EndTimestamp   int64  `json:"endTime"`
	Done           bool   `json:"done"`
	Reason         string `json:"reason"`
	Message        string `json:"message"`
}
