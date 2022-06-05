package schedule

import "errors"

const (
	DefaultJobListName    = "root"
	DefaultEventQueueName = "root"

	StateCreated = "Created"
	StateRunning = "Running"
	StateSuccess = "Success"
	StateFailed  = "Failed"

	ReasonServerError = "ServerError"
	ReasonClientError = "ClientError"
)

var (
	ErrorOverMaxScheduled    = errors.New("Max scheduled limit exceeded")
	ErrorOverMaxConcurrentcy = errors.New("Max sConcurrentcied limit exceeded")
	ErrorJobDoNotRuning      = errors.New("Cannot stop a job that is not running")
	ErrorJobNotFound         = errors.New("The job is not found")
	ErrorTimeout             = errors.New("Timeout")
)
