package api

import (
	"testing"
	"time"

	"github.com/madongming/scheduler/src/plugin"
)

func TestCreateJobList(t *testing.T) {
	type args struct {
		concurrency       int
		maxScheduledCount int
		maxJobCount       int
		maxHistory        int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create a Joblist and an EventQueue",
			args: args{
				concurrency:       100,
				maxScheduledCount: 100,
				maxHistory:        100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateJobList(tt.args.concurrency, tt.args.maxScheduledCount, tt.args.maxHistory); (err != nil) != tt.wantErr {
				t.Errorf("CreateJobList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if JobList == nil || EventQueue == nil {
				t.Errorf("JobList = %#v, EventQueue = %#v", JobList, EventQueue)
			}
		})
	}
}

func TestRunJob(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Run job 1",
			args: args{
				name: "job1",
			},
			wantErr: false,
		},
		{
			name: "Run job 2, can't runable",
			args: args{
				name: "job2",
			},
			wantErr: true,
		},
	}
	// 限制1个并发
	if err := CreateJobList(1, 10, 100); err != nil {
		t.Errorf("NewJobList(1, 10, 100) error = %v", err)
	}

	// 创建两个任务
	if err := CreateJob("job1", "任务1", plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3); err != nil {
		t.Errorf("CreateJob(job1, 任务1, plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3) error = %v", err)
	}
	if err := CreateJob("job2", "任务2", plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3); err != nil {
		t.Errorf("CreateJob(job2, 任务2, plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3) error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunJob(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RunJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getJob(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get a schedule job by name",
			args: args{
				name: "job1",
			},
			wantErr: false,
		},
	}

	// 准备数据
	if err := CreateJobList(10, 10, 100); err != nil {
		t.Errorf("NewJobList(10, 10, 100) error = %v", err)
	}

	if err := CreateJob("job1", "任务1", plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3); err != nil {
		t.Errorf("CreateJob(job1, 任务1, plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3) error = %v", err)
	}
	if err := CreateJob("job2", "任务2", plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3); err != nil {
		t.Errorf("CreateJob(job2, 任务2, plugin.SleepCmd1, 0, 5*time.Second, time.Second, 3) error = %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJob(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Jober == nil {
				t.Errorf("getJob() = %#v", got)
			}
		})
	}
}
