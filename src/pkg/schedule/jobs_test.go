package schedule

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/madongming/scheduler/src/pkg/store"
)

type JobShellScript struct {
}

func (js *JobShellScript) Run() error {
	time.Sleep(3 * time.Second)
	return nil
}

func (js *JobShellScript) Stop() error {
	return nil
}

func TestJobList_StartSchedule(t *testing.T) {
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "start job 1",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job1",
						Display:          "任务1",
						TypeName:         "ShellScript",
						ScheduleDuration: 10 * time.Second,
						Timeout:          5 * time.Second,
						RetryTimes:       3,
						RetryWait:        1 * time.Second,
					},
					Jober: &JobShellScript{},
				},
			},
			wantErr: false,
		},
		{
			name: "start job 2, can't runable",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job2",
						Display:          "任务2",
						TypeName:         "ShellScript",
						ScheduleDuration: 10 * time.Second,
						Timeout:          5 * time.Second,
						RetryTimes:       3,
						RetryWait:        1 * time.Second,
					},
					Jober: &JobShellScript{},
				},
			},
			wantErr: true,
		},
	}

	// 限制1个定时任务
	jl, err := NewJobList(10, 1, 100, true)
	if err != nil {
		t.Errorf("NewJobList(10, 1, 100, true) error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jl.StartSchedule(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("JobList.StartSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJobList_RunJob(t *testing.T) {
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Run job 1",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job1",
						Display:          "任务1",
						TypeName:         "ShellScript",
						ScheduleDuration: 10 * time.Second,
						Timeout:          5 * time.Second,
						RetryTimes:       3,
						RetryWait:        1 * time.Second,
					},
					Jober: &JobShellScript{},
				},
			},
			wantErr: false,
		},
	}

	jl, err := NewJobList(1, 10, 100, true)
	if err != nil {
		t.Errorf("NewJobList(1, 10, 100, true) error = %v", err)
	}

	for _, tt := range tests {
		go t.Run(tt.name, func(t *testing.T) {
			if err := jl.RunJob(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("JobList.RunJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_retry(t *testing.T) {
	type args struct {
		attempts uint8
		sleep    time.Duration
		fn       func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retry 3 time",
			args: args{
				attempts: 3,
				sleep:    time.Second,
				fn: func() error {
					fmt.Println("Exec 1")
					return errors.New("error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.attempts, tt.args.sleep, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
