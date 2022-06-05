package store

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var (
	testTime = time.Now()
)

type jobListStore struct {
	Key   string  `json:"key"`
	Value JobList `json:"value"`
}

type jobsStore []struct {
	Key   string `json:"key"`
	Value Job    `json:"value"`
}

type eventStore struct {
	Key   string     `json:"key"`
	Value EventQueue `json:"value"`
}

func TestCreateJobList(t *testing.T) {
	type args struct {
		Key   string  `json:"key"`
		Value JobList `json:"value"`
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantFileContent jobListStore
	}{
		{
			name: "Create Job List",
			args: args{
				Key: "root",
				Value: JobList{
					Concurrency:       100,
					MaxScheduledCount: 10,
					MaxHistory:        100,
					Jobs:              []string{},
				},
			},
			wantErr: false,
			wantFileContent: jobListStore{
				Key: "_jobs_default-root",
				Value: JobList{
					Concurrency:       100,
					MaxScheduledCount: 10,
					MaxHistory:        100,
					Jobs:              []string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateJobList(tt.args.Key, tt.args.Value); (err != nil) != tt.wantErr {
				t.Errorf("CreateJobList() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err := ioutil.ReadFile(JobListPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobListPath, err, tt.wantErr)
			}

			jobList := jobListStore{}
			err = json.Unmarshal(content, &jobList)
			if !reflect.DeepEqual(jobList, tt.wantFileContent) {
				t.Errorf("%s = %#v, want %#v", JobListPath, jobList, tt.wantFileContent)
			}
		})
	}
}

func TestAddJobListJobs(t *testing.T) {
	type args struct {
		name []string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantFileContent []string
	}{
		{
			name: "Add Job to job list",
			args: args{
				name: []string{"job1", "job2"},
			},
			wantFileContent: []string{"job1", "job2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddJobListJobs(tt.args.name...); (err != nil) != tt.wantErr {
				t.Errorf("AddJobListJobs() error = %v, wantErr %v", err, tt.wantErr)
			}
			content, err := ioutil.ReadFile(JobListPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobListPath, err, tt.wantErr)
			}

			jobList := jobListStore{}
			err = json.Unmarshal(content, &jobList)
			if !containArray(jobList.Value.Jobs, tt.wantFileContent) {
				t.Errorf("%s = %#v, want %#v", JobListPath, jobList, tt.wantFileContent)
			}
			DeleteJobListJobs(tt.args.name...)
		})
	}
}

func TestDeleteJobListJobs(t *testing.T) {
	type args struct {
		name []string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		wantFileContent []string
	}{
		{
			name: "Delete Job from job list",
			args: args{
				name: []string{"job1", "job2"},
			},
			wantFileContent: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddJobListJobs(tt.args.name...)

			if err := DeleteJobListJobs(tt.args.name...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteJobListJobs() error = %v, wantErr %v", err, tt.wantErr)
			}
			content, err := ioutil.ReadFile(JobListPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobListPath, err, tt.wantErr)
			}

			jobList := jobListStore{}
			err = json.Unmarshal(content, &jobList)
			if containArray(jobList.Value.Jobs, tt.wantFileContent) {
				t.Errorf("%s = %#v, want %#v", JobListPath, jobList, tt.wantFileContent)
			}
		})
	}
}

func TestAddJob(t *testing.T) {
	type args struct {
		value []Job
	}
	type twoValue struct {
		jobs        []Job
		jobListJobs []string
	}

	tests := []struct {
		name    string
		args    args
		want    twoValue
		wantErr bool
	}{
		{
			name: "Update a job",
			args: args{
				value: []Job{
					Job{
						Name:             "job1",
						Display:          "任务1",
						Type:             "ShellScript",
						ScheduleDuration: 3000000000,
						Timeout:          1000000000,
						RetryTimes:       3,
						RetryWait:        2000000000,
					},
				},
			},
			want: twoValue{
				jobs: []Job{
					Job{
						Name:             "job1",
						Display:          "任务1",
						Type:             "ShellScript",
						ScheduleDuration: 3000000000,
						Timeout:          1000000000,
						RetryTimes:       3,
						RetryWait:        2000000000,
					},
				},
				jobListJobs: []string{"job1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddJob(tt.args.value...); (err != nil) != tt.wantErr {
				t.Errorf("AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err := ioutil.ReadFile(JobListPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobListPath, err, tt.wantErr)
			}

			jobList := jobListStore{}
			err = json.Unmarshal(content, &jobList)

			content, err = ioutil.ReadFile(JobsPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobsPath, err, tt.wantErr)
			}

			jobs := jobsStore{}
			err = json.Unmarshal(content, &jobs)
			if err != nil {
				t.Errorf("json.Unmarshal(%v, %v) error = %v", content, jobs, err)
			}
			if !reflect.DeepEqual(jobList.Value.Jobs, tt.want.jobListJobs) {

				t.Errorf("%#v want %#v", jobList.Value.Jobs, tt.want.jobListJobs)
			}
			if !reflect.DeepEqual(jobs[0].Value, tt.want.jobs[0]) {
				t.Errorf("got:\n%#v \nwant:\n %#v", jobs[0].Value, tt.want.jobs[0])
			}
			DeleteJob(tt.args.value[0].Name)
		})
	}
}

func TestDeleteJob(t *testing.T) {
	type args struct {
		key []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete a job",
			args: args{
				key: []string{"job1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddJob(Job{
				Name:             "job1",
				Display:          "任务1",
				Type:             "ShellScript",
				ScheduleDuration: 3000000000,
				Timeout:          1000000000,
				RetryTimes:       3,
				RetryWait:        2000000000,
			})
			if err := DeleteJob(tt.args.key...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteJob() error = %v, wantErr %v", err, tt.wantErr)
			}
			content, err := ioutil.ReadFile(JobListPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobListPath, err, tt.wantErr)
			}

			jobList := jobListStore{}
			err = json.Unmarshal(content, &jobList)

			content, err = ioutil.ReadFile(JobsPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobsPath, err, tt.wantErr)
			}

			jobs := jobsStore{}
			err = json.Unmarshal(content, &jobs)
			if err != nil {
				t.Errorf("json.Unmarshal(%v, %v) error = %v", content, jobs, err)
			}
			if !reflect.DeepEqual(jobList.Value.Jobs, []string{}) {

				t.Errorf("%#v want %#v", jobList.Value.Jobs, []string{})
			}
			if !reflect.DeepEqual(jobs, jobsStore{}) {
				t.Errorf("got:\n%#v \nwant:\n %#v", jobs, jobsStore{})
			}
		})
	}
}

func TestUpdateJob(t *testing.T) {
	type args struct {
		key   string
		value Job
	}
	tests := []struct {
		name    string
		args    args
		want    Job
		wantErr bool
	}{
		{
			name: "Update a job",
			args: args{
				key: "job1",
				value: Job{
					Name:             "job1",
					Display:          "任务1",
					Type:             "ShellScript",
					ScheduleDuration: 4000000000,
					Timeout:          1000000000,
					RetryTimes:       3,
					RetryWait:        2000000000,
				},
			},
			want: Job{
				Name:             "job1",
				Display:          "任务1",
				Type:             "ShellScript",
				ScheduleDuration: 4000000000,
				Timeout:          1000000000,
				RetryTimes:       3,
				RetryWait:        2000000000,
			},
			wantErr: false,
		},
	}
	AddJob(Job{
		Name:             "job1",
		Display:          "任务1",
		Type:             "ShellScript",
		ScheduleDuration: 3000000000,
		Timeout:          1000000000,
		RetryTimes:       3,
		RetryWait:        2000000000,
	})
	defer DeleteJob("job1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UpdateJob(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			content, err := ioutil.ReadFile(JobsPath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", JobsPath, err, tt.wantErr)
			}

			jobs := jobsStore{}
			err = json.Unmarshal(content, &jobs)
			if err != nil {
				t.Errorf("json.Unmarshal(%v, %v) error = %v", content, jobs, err)
			}
			if !reflect.DeepEqual(jobs[0].Value, tt.want) {
				t.Errorf("got:\n%#v \nwant:\n%#v", jobs[0].Value, tt.want)
			}
		})
	}
}

func TestGetJobs(t *testing.T) {
	tests := []struct {
		name    string
		want    []Job
		wantErr bool
	}{
		{
			name: "Get a Job list",
			want: []Job{
				Job{
					Name:             "job1",
					Display:          "任务1",
					Type:             "ShellScript",
					ScheduleDuration: 3000000000,
					Timeout:          1000000000,
					RetryTimes:       3,
					RetryWait:        2000000000,
				},
			},
			wantErr: false,
		},
	}
	AddJob(Job{
		Name:             "job1",
		Display:          "任务1",
		Type:             "ShellScript",
		ScheduleDuration: 3000000000,
		Timeout:          1000000000,
		RetryTimes:       3,
		RetryWait:        2000000000,
	})
	defer DeleteJob("job1")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetJobs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Job
		wantErr bool
	}{
		{
			name: "Get a job by name",
			args: args{
				name: "job1",
			},
			want: Job{
				Name:             "job1",
				Display:          "任务1",
				Type:             "ShellScript",
				ScheduleDuration: 3000000000,
				Timeout:          1000000000,
				RetryTimes:       3,
				RetryWait:        2000000000,
			},
			wantErr: false,
		},
		{
			name: "Get a not exsit job by name",
			args: args{
				name: "job2",
			},
			wantErr: true,
		},
	}
	AddJob(Job{
		Name:             "job1",
		Display:          "任务1",
		Type:             "ShellScript",
		ScheduleDuration: 3000000000,
		Timeout:          1000000000,
		RetryTimes:       3,
		RetryWait:        2000000000,
	})
	defer DeleteJob("job1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetJob(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateEventQueue(t *testing.T) {
	type args struct {
		key   string
		value EventQueue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    eventStore
	}{
		{
			name: "Create a event queue",
			args: args{
				key: "root",
				value: EventQueue{
					MaxStorage: 100,
					MaxHistory: 100,
				},
			},
			wantErr: false,
			want: eventStore{
				Key: "_event_default-root",
				Value: EventQueue{
					MaxStorage: 100,
					MaxHistory: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateEventQueue(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("CreateEventQueue() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err := ioutil.ReadFile(EventQueuePath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", EventQueuePath, err, tt.wantErr)
			}
			event := eventStore{}
			err = json.Unmarshal(content, &event)
			if err != nil {
				t.Errorf("json.Unmarshal(%v, %v) error = %v", content, event, err)
			}

			if !reflect.DeepEqual(event, tt.want) {
				t.Errorf("got:\n%#v \nwant:\n%#v", event, tt.want)
			}
		})
	}
}

func TestAddHistory(t *testing.T) {
	type args struct {
		value []JobInstance
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    eventStore
	}{
		{
			name: "Add s history",
			args: args{
				value: []JobInstance{
					JobInstance{
						Name:    "job1",
						Display: "任务1",
						Type:    "ShellScript",
						State:   "Success",
						Done:    true,
						Reason:  "",
						Message: "Hello, world",
					},
				},
			},
			wantErr: false,
			want: eventStore{
				Key: "_event_default-root",
				Value: EventQueue{
					MaxStorage: 100,
					MaxHistory: 100,
					History: []JobInstance{
						JobInstance{
							Name:    "job1",
							Display: "任务1",
							Type:    "ShellScript",
							State:   "Success",
							Done:    true,
							Reason:  "",
							Message: "Hello, world",
						},
					},
				},
			},
		},
	}
	CreateEventQueue("root", EventQueue{
		MaxStorage: 100,
		MaxHistory: 100})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddHistory(tt.args.value...); (err != nil) != tt.wantErr {
				t.Errorf("AddHistory() error = %v, wantErr %v", err, tt.wantErr)
			}

			content, err := ioutil.ReadFile(EventQueuePath)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%s) error = %v, wantErr %v", EventQueuePath, err, tt.wantErr)
			}
			event := eventStore{}
			err = json.Unmarshal(content, &event)
			if err != nil {
				t.Errorf("json.Unmarshal(%v, %v) error = %v", content, event, err)
			}

			if !reflect.DeepEqual(event, tt.want) {
				t.Errorf("got:\n%#v \nwant:\n%#v", event, tt.want)
			}
		})
	}
}

func TestGetHistory(t *testing.T) {
	tests := []struct {
		name    string
		want    []JobInstance
		wantErr bool
	}{
		{
			name: "Get the history",
			want: []JobInstance{
				JobInstance{
					Name:    "job1",
					Display: "任务1",
					Type:    "ShellScript",
					State:   "Success",
					Done:    true,
					Reason:  "",
					Message: "Hello, world",
				},
			},
			wantErr: false,
		},
	}
	CreateEventQueue("root", EventQueue{
		MaxStorage: 100,
		MaxHistory: 100})
	AddHistory(JobInstance{
		Name:    "job1",
		Display: "任务1",
		Type:    "ShellScript",
		State:   "Success",
		Done:    true,
		Reason:  "",
		Message: "Hello, world",
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHistory()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isInArray(t *testing.T) {
	type args struct {
		str   string
		array []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "The string in the array",
			args: args{
				str:   "abc",
				array: []string{"abc", "cde"},
			},
			want: true,
		},
		{
			name: "The string is not in the array",
			args: args{
				str:   "def",
				array: []string{"abc", "cde"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInArray(tt.args.str, tt.args.array); got != tt.want {
				t.Errorf("isInArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

// util func
// containArray two array, a1 contain a2. return true
func containArray(a1, a2 []string) bool {
	if len(a1) == 0 {
		return false
	}

	for _, s := range a2 {
		if !containStr(a1, s) {
			return false
		}
	}
	return true
}

// containStr if the string str in array, return true
func containStr(array []string, str string) bool {
	for _, s := range array {
		if str == s {
			return true
		}
	}
	return false
}
