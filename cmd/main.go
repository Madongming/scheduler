package main

import (
	"fmt"
	"time"

	"github.com/madongming/scheduler/src/pkg/api"
)

const (
	concurrency       = 10
	maxScheduledCount = 10
	maxHistory        = 100
)

func main() {
	// 初始化
	fmt.Println("创建调度器...")
	fmt.Printf("可同时执行 %d 个任务，后台可以管理定时任务 %d 个，最多存储 %d 个执行记录...\n",
		concurrency, maxScheduledCount, maxHistory)
	if err := api.CreateJobList(concurrency, maxScheduledCount, maxHistory); err != nil {
		panic(err)
	}
	fmt.Println("创建完成")
	fmt.Println("=================")
	fmt.Println("")

	// 任务的增｜删｜改｜查
	fmt.Println("创建任务:")
	fmt.Println("任务1:")
	fmt.Println("名称：job1，显示名：任务1，类型名：SleepCmd1")
	fmt.Println("定时执行：否，执行超时：5s，重试次数：3次，重试间隔：1s")
	if err := api.CreateJob("job1", "任务1", "SleepCmd1", 0, 5*time.Second, time.Second, 3); err != nil {
		panic(err)
	}
	fmt.Println("创建完成")
	fmt.Println("-----------------")
	fmt.Println()

	fmt.Println("创建任务:")
	fmt.Println("任务2:")
	fmt.Println("名称：job2，显示名：任务2，类型名：SleepCmd1")
	fmt.Println("定时执行：否，执行超时：5s，重试次数：3次，重试间隔：1s")
	if err := api.CreateJob("job2", "任务2", "SleepCmd2", 0, 5*time.Second, time.Second, 3); err != nil {
		panic(err)
	}
	fmt.Println("创建完成")
	fmt.Println("-----------------")
	fmt.Println()

	fmt.Println("查看任务列表")
	jobList, err := api.GetJobList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", jobList)
	fmt.Println("-----------------")
	fmt.Println()

	fmt.Println("删除任务2")
	if err := api.DeleteJobList("job2"); err != nil {
		panic(err)
	}
	fmt.Println("删除任务2完成，查看任务列表")
	jobList, err = api.GetJobList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", jobList)
	fmt.Println("-----------------")
	fmt.Println()

	fmt.Println("修改任务1")
	scheduleJob, err := api.GetScheduleJob("job1")
	if err != nil {
		panic(err)
	}
	scheduleJob.Display = "新job1"
	if _, err := api.UpdateJobList(scheduleJob); err != nil {
		panic(err)
	}
	fmt.Println("更新任务1完成，查看任务")
	job, err := api.GetJob("job1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", job)
	fmt.Println("查看任务列表")
	jobList, err = api.GetJobList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", jobList)
	fmt.Println("=================")
	fmt.Println()
	fmt.Println()

	// 任务运行
	fmt.Println("运行任务1")
	if err := api.RunJob("job1"); err != nil {
		panic(err)
	}
	fmt.Println("查看正在运行的任务：")
	running, err := api.RunningJob()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", running)
	time.Sleep(5 * time.Second)
	fmt.Println("运行完成")
	fmt.Println("查看正在运行的任务：")
	running, err = api.RunningJob()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", running)
	fmt.Println("-----------------")
	fmt.Println()

	fmt.Println("创建定时执行的任务:")
	fmt.Println("任务2:")
	fmt.Println("名称：job2，显示名：任务2，类型名：SleepCmd2")
	fmt.Println("定时执行：每5s，执行超时：5s，重试次数：3次，重试间隔：1s")
	if err := api.CreateJob("job2", "任务2", "SleepCmd2", 5*time.Second, 5*time.Second, time.Second, 3); err != nil {
		panic(err)
	}
	fmt.Println("创建完成,次任务会每隔5s执行一次，等待20s")
	fmt.Println("-----------------")
	fmt.Println()

	time.Sleep(5 * time.Second)
	fmt.Println("同时执行任务1，并行的体现")
	if err := api.RunJob("job1"); err != nil {
		panic(err)
	}

	time.Sleep(15 * time.Second)
	fmt.Println("=================")
	fmt.Println()

	fmt.Println("查看执行历史")
	history, err := api.JobHistory()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", history)
}
