package schedule

import (
	"reflect"
	"testing"
	"time"

	"github.com/madongming/scheduler/src/pkg/store"
)

func TestEventQueue_Push(t *testing.T) {
	type args struct {
		job *Job
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Push a job to the queue",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job1",
						Display:          "任务1",
						Type:             "ShellScript",
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
			name: "Push the second job to the queue",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job2",
						Display:          "任务2",
						Type:             "ShellScript",
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
			name: "Push the third job to the queue, failed",
			args: args{
				job: &Job{
					Job: store.Job{
						Name:             "job3",
						Display:          "任务3",
						Type:             "ShellScript",
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

	// 限制2个job
	eq, err := NewEventQueue(2, 100, true)
	if err != nil {
		t.Errorf("NewEventQueue(2, 100, true) error = %v", err)
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := eq.Push(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("EventQueue.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(eq.events) != i+1 && i < 2 {
				t.Errorf("Push %d failed, %v", i+1, eq.events)
			}
		})
	}
}

func TestEventQueue_Pop2History(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Pop a Job to the history",
			args: args{
				name: "job1",
			},
			wantErr: false,
		},
		{
			name: "Pop second Job to the history",
			args: args{
				name: "job2",
			},
			wantErr: false,
		},
	}
	// 准备数据
	// 限制2个job, history限制1个
	eq, err := NewEventQueue(2, 1, true)
	if err != nil {
		t.Errorf("NewEventQueue(2, 1, true) error = %v", err)
	}
	preData := []Job{
		Job{
			Job: store.Job{
				Name:             "job1",
				Display:          "任务1",
				Type:             "ShellScript",
				ScheduleDuration: 10 * time.Second,
				Timeout:          5 * time.Second,
				RetryTimes:       3,
				RetryWait:        1 * time.Second,
			},
			Jober: &JobShellScript{},
		},
		Job{
			Job: store.Job{
				Name:             "job2",
				Display:          "任务2",
				Type:             "ShellScript",
				ScheduleDuration: 10 * time.Second,
				Timeout:          5 * time.Second,
				RetryTimes:       3,
				RetryWait:        1 * time.Second,
			},
			Jober: &JobShellScript{},
		},
	}
	for i := range preData {
		eq.Push(&preData[i])
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := eq.Pop2History(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("EventQueue.Pop2History() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(eq.History) != 1 {
				t.Errorf("%#v", eq.History)
			}
		})
	}
}

func Test_deleteInstanceByName(t *testing.T) {
	type args struct {
		name   string
		events []store.JobInstance
	}
	tests := []struct {
		name    string
		args    args
		want    store.JobInstance
		want1   []store.JobInstance
		wantErr bool
	}{
		{
			name: "Delete by an exsit name",
			args: args{
				name: "job1",
				events: []store.JobInstance{
					store.JobInstance{
						Name: "job1",
					},
					store.JobInstance{
						Name: "job2",
					},
				},
			},
			want: store.JobInstance{Name: "job1"},
			want1: []store.JobInstance{
				store.JobInstance{
					Name: "job2",
				},
			},
			wantErr: false,
		},
		{
			name: "Delete by a not exsit name",
			args: args{
				name: "job3",
				events: []store.JobInstance{
					store.JobInstance{
						Name: "job1",
					},
					store.JobInstance{
						Name: "job2",
					},
				},
			},
			want: store.JobInstance{},
			want1: []store.JobInstance{
				store.JobInstance{
					Name: "job1",
				},
				store.JobInstance{
					Name: "job2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := deleteInstanceByName(tt.args.name, tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteInstanceByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteInstanceByName() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("deleteInstanceByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
