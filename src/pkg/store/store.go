package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	JobListPrefix    = "_jobs_default"
	EventQueuePrefix = "_event_default"

	JobListPath    = "./mock/joblist.json"
	JobsPath       = "./mock/jobs.json"
	EventQueuePath = "./mock/eventQueue.json"
)

// CreateJobList Create JobList store
func CreateJobList(key string, value JobList) error {
	content, err := json.Marshal(
		struct {
			Key   string  `json:"key"`
			Value JobList `json:"value"`
		}{
			Key:   fmt.Sprintf("%s-%s", JobListPrefix, key),
			Value: value,
		})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(JobListPath, content, 0644)
}

// AddJobListJobs add the Job names into JobList
func AddJobListJobs(name ...string) error {
	content, err := ioutil.ReadFile(JobListPath)
	if err != nil {
		return err
	}

	jobList := struct {
		Key   string  `json:"key"`
		Value JobList `json:"value"`
	}{}
	err = json.Unmarshal(content, &jobList)

	jobList.Value.Jobs = append(jobList.Value.Jobs, name...)
	content, err = json.Marshal(&jobList)
	return ioutil.WriteFile(JobListPath, content, 0644)
}

// DeleteJobListJobs delete the Job names from JobList
func DeleteJobListJobs(name ...string) error {
	content, err := ioutil.ReadFile(JobListPath)
	if err != nil {
		return err
	}

	jobList := struct {
		Key   string  `json:"key"`
		Value JobList `json:"value"`
	}{}
	err = json.Unmarshal(content, &jobList)

	newJobs := make([]string, 0, len(jobList.Value.Jobs))
	for i := range jobList.Value.Jobs {
		if isInArray(jobList.Value.Jobs[i], name) {
			continue
		}
		newJobs = append(newJobs, jobList.Value.Jobs[i])
	}
	jobList.Value.Jobs = newJobs

	content, err = json.Marshal(&jobList)
	return ioutil.WriteFile(JobListPath, content, 0644)
}

// AddJob Add a job to store, and add the name to the JobList
func AddJob(value ...Job) error {
	data := make([]struct {
		Key   string `json:"key"`
		Value Job    `json:"value"`
	}, len(value))
	names := make([]string, len(value))

	for i, v := range value {
		data[i].Key = v.Name
		data[i].Value = v
		names[i] = v.Name
	}

	if err := AddJobListJobs(names...); err != nil {
		return err
	}

	var (
		content []byte
		err     error
	)
	if _, err = os.Stat(JobsPath); os.IsNotExist(err) {
		content, err = json.Marshal(data)
		if err != nil {
			return err
		}
	} else {
		content, err = ioutil.ReadFile(JobsPath)
		if err != nil {
			return err
		}
		var oldData []struct {
			Key   string `json:"key"`
			Value Job    `json:"value"`
		}

		err = json.Unmarshal(content, &oldData)
		if err != nil {
			return err
		}

		oldData = append(oldData, data...)
		content, err = json.Marshal(oldData)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(JobsPath, content, 0644)
}

// DeleteJob Delete some job from store, and delete the names from the JobList
func DeleteJob(key ...string) error {
	if err := DeleteJobListJobs(key...); err != nil {
		return err
	}
	content, err := ioutil.ReadFile(JobsPath)
	if err != nil {
		return err
	}
	var data []struct {
		Key   string `json:"key"`
		Value Job    `json:"value"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	newJobs := make([]struct {
		Key   string `json:"key"`
		Value Job    `json:"value"`
	}, len(data))
	index := 0
	for i := range data {
		if isInArray(data[i].Key, key) {
			continue
		}
		newJobs[index].Key = data[i].Key
		newJobs[index].Value = data[i].Value
		index++
	}
	data = newJobs[:index]

	content, err = json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(JobsPath, content, 0644)
}

func UpdateJob(key string, value Job) (Job, error) {
	content, err := ioutil.ReadFile(JobsPath)
	if err != nil {
		return Job{}, err
	}
	var data []struct {
		Key   string `json:"key"`
		Value Job    `json:"value"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return Job{}, err
	}
	for i := range data {
		if data[i].Key == key {
			data[i].Value = value
			break
		}
	}

	content, err = json.Marshal(data)
	if err != nil {
		return Job{}, err
	}

	return value, ioutil.WriteFile(JobsPath, content, 0644)
}
func GetJobs() ([]Job, error) {
	content, err := ioutil.ReadFile(JobsPath)
	if err != nil {
		return []Job{}, err
	}
	var data []struct {
		Key   string `json:"key"`
		Value Job    `json:"value"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return []Job{}, err
	}
	re := make([]Job, len(data))
	for i := range data {
		re[i] = data[i].Value
	}

	return re, nil
}

// Events
func CreateEventQueue(key string, value EventQueue) error {
	content, err := json.Marshal(
		struct {
			Key   string     `json:"key"`
			Value EventQueue `json:"value"`
		}{
			Key:   fmt.Sprintf("%s-%s", EventQueuePrefix, key),
			Value: value,
		})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(EventQueuePath, content, 0644)

}

func AddHistory(value ...JobInstance) error {
	content, err := ioutil.ReadFile(EventQueuePath)
	if err != nil {
		return err
	}
	var data struct {
		Key   string     `json:"key"`
		Value EventQueue `json:"value"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	data.Value.History = append(data.Value.History, value...)
	if len(data.Value.History) > data.Value.MaxHistory {
		data.Value.History = data.Value.History[len(data.Value.History)-data.Value.MaxHistory:]
	}
	content, err = json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(EventQueuePath, content, 0644)
}
func GetHistory() ([]JobInstance, error) {
	content, err := ioutil.ReadFile(EventQueuePath)
	if err != nil {
		return []JobInstance{}, err
	}
	var data struct {
		Key   string     `json:"key"`
		Value EventQueue `json:"value"`
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return []JobInstance{}, err
	}

	return data.Value.History, nil
}

// util functions
func isInArray(str string, array []string) bool {
	for _, s := range array {
		if str == s {
			return true
		}
	}
	return false
}
