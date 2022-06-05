package api

import (
	"time"

	"github.com/madongming/scheduler/src/pkg/schedule"
	"github.com/madongming/scheduler/src/pkg/store"
	"github.com/madongming/scheduler/src/plugin"
)

var (
	// 真实使用的时候可以不设置为全局变量，可以设置多个列表
	JobList    *schedule.JobList
	EventQueue *schedule.EventQueue
)

func CreateJobList(concurrency, maxScheduledCount, maxHistory int) error {
	var err error
	JobList, err = schedule.NewJobList(concurrency, maxScheduledCount, maxHistory, true)
	if err != nil {
		return err
	}

	EventQueue, err = schedule.NewEventQueue(concurrency, maxHistory, true)
	if err != nil {
		return err
	}
	return nil
}

func CreateJob(name, display, typeName string, scheduleDuration, timeout, retryWait time.Duration, retryTimes uint8) error {
	var jober schedule.Jober
	switch typeName {
	case plugin.SleepCmd1:
		jober = plugin.NewSleep1()
	case plugin.SleepCmd2:
		jober = plugin.NewSleep2()
	default:
		return ErrorPluginNotFound
	}
	job, err := schedule.NewJob(name, display, typeName, scheduleDuration, timeout, retryWait, retryTimes, jober)
	if err != nil {
		return err
	}

	if job.ScheduleDuration != time.Duration(0) {
		// 一个需要定时执行的job
		if err := JobList.StartSchedule(&job); err != nil {
			return err
		}
	}

	return JobList.Add(job)
}

func DeleteJobList(name string) error {
	return JobList.Delete(name)
}

func UpdateJobList(job schedule.Job) (schedule.Job, error) {
	return JobList.Update(job)
}

func GetJobList() ([]store.Job, error) {
	return JobList.Get()
}

func GetJob(name string) (store.Job, error) {
	return JobList.GetJob(name)
}

func RunJob(name string) error {
	job, err := getJob(name)
	if err != nil {
		return err
	}
	if err := EventQueue.Push(&job); err != nil {
		return err
	}

	go runJob(&job, EventQueue)

	return nil
}

func StopJob(name string) error {
	job, err := getJob(name)
	if err != nil {
		return err
	}
	return JobList.StopJob(&job)
}

func RunningJob() ([]store.JobInstance, error) {
	return schedule.RunningJob(EventQueue)
}

func JobHistory() ([]store.JobInstance, error) {
	return schedule.JobHistory(EventQueue)
}

func getJob(name string) (schedule.Job, error) {
	storeJob, err := JobList.GetJob(name)
	if err != nil {
		return schedule.Job{}, err
	}
	var job schedule.Job

	switch storeJob.TypeName {
	case plugin.SleepCmd1:
		job, err = schedule.NewJob(
			storeJob.Name,
			storeJob.Display,
			storeJob.TypeName,
			storeJob.ScheduleDuration,
			storeJob.Timeout,
			storeJob.RetryWait,
			storeJob.RetryTimes,
			plugin.NewSleep1())
		if err != nil {
			return schedule.Job{}, err
		}
	case plugin.SleepCmd2:
		job, err = schedule.NewJob(
			storeJob.Name,
			storeJob.Display,
			storeJob.TypeName,
			storeJob.ScheduleDuration,
			storeJob.Timeout,
			storeJob.RetryWait,
			storeJob.RetryTimes,
			plugin.NewSleep2())
		if err != nil {
			return schedule.Job{}, err
		}
	default:
		return schedule.Job{}, ErrorPluginNotFound
	}
	return job, nil
}

func runJob(job *schedule.Job, eventQueue *schedule.EventQueue) error {
	defer eventQueue.Pop2History(job.Name)
	return JobList.RunJob(job)
}
