package schedule

import (
	"os"
	"time"

	"github.com/madongming/scheduler/src/pkg/store"
)

type JobList struct {
	store.JobList

	chConcurrentcy chan struct{}
	chScheduled    chan struct{}
}

type Job struct {
	store.Job

	Jober Jober

	ticker *time.Ticker
}

func NewJobList(concurrency, maxScheduledCount, maxHistory int, force ...bool) (*JobList, error) {
	jobList := JobList{
		JobList: store.JobList{
			Concurrency:       concurrency,
			MaxScheduledCount: maxScheduledCount,
			MaxHistory:        maxHistory,
			Jobs:              []string{},
		},

		chConcurrentcy: make(chan struct{}, concurrency),
		chScheduled:    make(chan struct{}, maxScheduledCount),
	}
	if _, err := os.Stat(store.JobsPath); os.IsNotExist(err) ||
		(force != nil &&
			len(force) > 0 &&
			force[0]) {
		if err = store.CreateJobList(DefaultJobListName, jobList.JobList); err != nil {
			return nil, err
		}
	}
	return &jobList, nil
}

func NewJob(name, display, typeName string, scheduleDuration, timeout, retryWait time.Duration, retryTimes uint8, jober Jober) (Job, error) {
	return Job{
		Job: store.Job{
			Name:             name,
			Display:          display,
			Type:             typeName,
			ScheduleDuration: scheduleDuration,
			Timeout:          timeout,
			RetryTimes:       retryTimes,
			RetryWait:        retryWait,
		},
		Jober: jober,
	}, nil
}

func (_ *JobList) Add(job Job) error {
	return store.AddJob(job.Job)
}

func (_ *JobList) Delete(name string) error {
	return store.DeleteJob(name)
}

func (_ *JobList) Update(job Job) (Job, error) {
	_, err := store.UpdateJob(job.Name, job.Job)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func (_ *JobList) Get() ([]store.Job, error) {
	return store.GetJobs()
}

func (_ *JobList) GetJob(name string) (store.Job, error) {
	return store.GetJob(name)
}

func (jl *JobList) StartSchedule(job *Job) error {
	select {
	case jl.chScheduled <- struct{}{}:
		go func(j *Job) {
			defer func() { <-jl.chScheduled }()
			j.ticker = time.NewTicker(j.ScheduleDuration)
			for _ = range j.ticker.C {
				jl.RunJob(j)
			}
		}(job)
	default:
		return ErrorOverMaxScheduled
	}
	return nil
}

func (jl *JobList) StopSchedule(job *Job) error {
	if job.ticker == nil {
		return ErrorJobDoNotRuning
	}
	job.ticker.Stop()

	return nil
}

func (jl *JobList) RunJob(job *Job) error {
	select {
	case jl.chConcurrentcy <- struct{}{}:
		go func() {
			defer func() { <-jl.chConcurrentcy }()

			timer := time.NewTimer(job.Timeout)
			done := make(chan error, 1)

			go func(done chan error) {
				done <- retry(job.RetryTimes, job.RetryWait, job.Jober.Run)
			}(done)

			select {
			case <-timer.C:
				return
			case <-done:
				return
			}
		}()
	default:
		return ErrorOverMaxConcurrentcy
	}
	return nil
}

func (jl *JobList) StopJob(job *Job) error {
	return job.Jober.Stop()
}

func retry(attempts uint8, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, sleep, fn)
		}
		return err
	}
	return nil
}
