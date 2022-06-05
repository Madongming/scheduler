package schedule

import (
	"os"
	"sync"
	"time"

	"github.com/madongming/scheduler/src/pkg/store"
)

type EventQueue struct {
	store.EventQueue

	events []store.JobInstance

	mu *sync.Mutex // Make queue operations safe
}

//
func NewEventQueue(maxStorage, maxHistory int, force ...bool) (*EventQueue, error) {
	eventQueue := EventQueue{
		EventQueue: store.EventQueue{
			MaxStorage: maxStorage,
			MaxHistory: maxHistory,

			History: make([]store.JobInstance, 0, maxHistory),
		},

		events: make([]store.JobInstance, 0, maxStorage),
		mu:     &sync.Mutex{},
	}
	if _, err := os.Stat(store.EventQueuePath); os.IsNotExist(err) ||
		(force != nil &&
			len(force) > 0 &&
			force[0]) {
		if err = store.CreateEventQueue(DefaultEventQueueName, eventQueue.EventQueue); err != nil {
			return nil, err
		}
	}
	return &eventQueue, nil
}

func (eq *EventQueue) Push(job *Job) error {
	jobInstance, err := store.NewJobInstance(&job.Job)
	jobInstance.State = StateRunning
	jobInstance.StartTime = time.Now()
	if err != nil {
		return err
	}

	eq.mu.Lock()
	defer eq.mu.Unlock()

	if len(eq.events) == eq.MaxStorage {
		return ErrorOverMaxConcurrentcy
	}
	eq.events = append(eq.events, jobInstance)

	return nil
}

func (eq *EventQueue) Pop2History(name string) error {
	eq.mu.Lock()
	defer eq.mu.Unlock()

	var (
		jobInstance store.JobInstance
		err         error
	)

	jobInstance, eq.events, err = deleteInstanceByName(name, eq.events)
	if err != nil {
		return err
	}

	jobInstance.State = StateSuccess
	jobInstance.Done = true
	jobInstance.EndTime = time.Now()

	err = store.AddHistory(jobInstance)
	if err != nil {
		return err
	}
	eq.EventQueue.History, err = store.GetHistory()
	if err != nil {
		return err
	}
	return nil
}

func deleteInstanceByName(name string, events []store.JobInstance) (store.JobInstance, []store.JobInstance, error) {
	var (
		jobInstance store.JobInstance
		found       bool
	)
	for i := range events {
		if events[i].Name == name {
			found = true
			jobInstance = events[i]
			if i == len(events)-1 {
				// 是最后一个元素
				break
			}
			for j, e := range events[i+1:] {
				events[i+j] = e
			}
			break
		}
	}
	if found {
		events = events[:len(events)-1]
		return jobInstance, events, nil
	}
	return jobInstance, events, ErrorJobNotFound
}

func RunningJob(eq *EventQueue) ([]store.JobInstance, error) {
	return eq.events, nil
}

func JobHistory(_ *EventQueue) ([]store.JobInstance, error) {
	return store.GetHistory()
}
