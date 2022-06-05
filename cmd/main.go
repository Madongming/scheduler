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
	fmt.Println("创建调度器...")
	fmt.Printf("可同时执行 %d 个任务，后台可以管理定时任务 %d 个，最多存储 %d 个执行记录...\n",
		concurrency, maxScheduledCount, maxHistory)
	if err := api.CreateJobList(concurrency, maxScheduledCount, maxHistory); err != nil {
		panic(err)
	}
	fmt.Println("创建完成")
	fmt.Println("=================")
	fmt.Println("")

	fmt.Println("创建任务:")
	fmt.Println("任务1:")
	fmt.Println("名称：job1，显示名：任务1，类型名：SleepCmd1")
	fmt.Println("定时执行：否，执行超时：5s，重试次数：3次，重试间隔：1s")
	if err := api.CreateJob("job1", "任务1", "SleepCmd1", 0, 5*time.Second, time.Second, 3); err != nil {
		panic(err)
	}
	fmt.Println("创建完成")
	fmt.Println("-----------------")
	fmt.Println("创建任务:")
	fmt.Println("任务2:")
	fmt.Println("名称：job2，显示名：任务2，类型名：SleepCmd2")
	fmt.Println("定时执行：每10s，执行超时：5s，重试次数：3次，重试间隔：1s")
	if err := api.CreateJob("job2", "任务2", "SleepCmd2", 10*time.Second, 5*time.Second, time.Second, 3); err != nil {
		panic(err)
	}
	fmt.Println("创建完成,次任务会每隔10s执行一次，等待30s")
	time.Sleep(30 * time.Second)
	fmt.Println("-----------------")
	fmt.Println("运行任务1")
	if err := api.RunJob("job1"); err != nil {
		panic(err)
	}
	fmt.Println("等待30s 观察状态")
	time.Sleep(30 * time.Second)
	fmt.Println("运行完成")
}
