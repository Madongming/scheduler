# 数据

## 数据结构
### 任务
``` golang
type JobList struct {
	Concurrency int `json:"concurrency"`
	MaxScheduledCount int `json:"maxScheduledCount"`
	MaxHistory int `json:"maxHistory"`
	
	Jobs []Job `json:"jobs"`
}
```

### 每个任务
``` golang
type Job struct {
    Name string `json:"name"`
	Display string `json:"display"`
	Type string `json:"type"`
	ScheduleDuration time.Duration `json:"scheduleDuration"`
	Timeout time.Duration `json:"timeout"`
	RetryTimes uint8 `json:"retryTimes"`
	RetryWait time.Duration `json:"retryWait"` // 默认等于超时时间
	Jober Jober `json:"spec"`
}
```
每个任务的实例
``` golang
type JobInstance struct {
    Name string `json:"name"`
	Display string `json:"display"`
	Type string `json:"type"`
    State string `json:"state"`
	StartTimestamp int64 `json:"startTime"`
	EndTimestamp int64 `json:"endTime"`
	Done bool `json:"done"`
	Reason string `json:"reason"`
	Message string `json:"message"`
}
```

## 事件队列
``` golang
type EventQueue struct {
	MaxStorage int `json:"maxStorage"`
	MaxHistory int `json:"maxHistory"`
	
    events []*JobInstance
	history[]*JobInstance
}
```

## 持久化
- JobList
``` json
{
    "key": "default",
	"value": {
        "concurrency": <number>,
	    "maxScheduledCount": <number>,
	    "maxHistory": <number>,
	    "jobs": [
	        "jobs1",
	    	"job2"
        ]
	}
}
```
- Job
``` json
[
    {
	    "key": "job1",
		"value": {
            "name": "job1",
	        "display": "任务1",
            "type": "ShellScript",
            "scheduleDuration": 3000000000,
            "timeout": 10000000000,
            "retryTimes": 3,
	        "retryWait": 3000000000
		}
	},
	{
	    "key": "job2",
		"value": {
	        "name": "job2",
	        "display": "任务2",
            "type": "ShellScript",
            "scheduleDuration": 3000000000,
            "timeout": 10000000000,
            "retryTimes": 3,
	        "retryWait": 3000000000
		}
	}
]
```
- EventQueue
```json
{
    "key": "default": {
        "maxStorage": <number>,
        "maxHistory": <number>,
		"history": [
		    "job1",
			"job2"
		]
	}
}
    
```
- Event
```json
[
    {
        "name": "job1",
        "display": "任务1",
        "type": "ShellScript",
        "state": "Sueccess",
        "startTime": 1523443233,
        "endTime": 1523443433,
        "done": true,
        "reason": "",
        "message":"定时执行"
	},
	{
        "name": "job2",
        "display": "任务2",
        "type": "ShellScript",
        "state": "Failed",
        "startTime": 1523443233,
        "endTime": 1523443433,
        "done": true,
        "reason": "ClinetError",
        "message":"重试第2次，Command not found"
	},
]


## 数据行为（methods）

### JobList
- Add(job Job) error // 增加一个Job
- Delete(name string) error // 删除一个Job
- Update(job Job) (Job, error) // 更新一个Job
- Get() (JobList, error) // 列出所有Job
- GetJob(name string) (Job, error) // 获取指定Job
- StartSchedule(job Job) error // 如果是需要调度执行的Job，开始这个调度
- StopSchedule(job Job) error // 如果是需要调度执行的Job，停止这个调度

### Job
- Run(ctx context) error // 运行Job
- Stop(ctx context) error // 停止Job

### EventQueue
- Push(event Event) error // 加入一个Job
- Pop() (event Event, error) // 取出一个Job

### Event
- Process() error // 处理一个事件

# 行为

## 行为抽象（interface）
### 任务接口
``` golang
type Jober interface {
    Run(ctx context) error
	Stop(ctx context) error
}
```

## 应用层
### 工厂函数
- JobList New(concurrency, maxScheduledCount int) (*JobList, error)
- Job New(name, display, type string, scheduleDuration, timeout time.Duration, jober Jober) (*Job, error)
- EventQueue New(maxStorage int) (*EventQueue, error)
- Event New(name, type string) (*Event, error)

### 业务函数
#### JobList处理
- AddJobList(*Job) error // 添加一个Job，如果是定期执行的，启动定期执行的计时器
- DeleteJobList(name string) error // 删除一个Job，如果是定期执行的，关闭定期执行的计时器
- UpdateJobList(job Job) (Job, error) // 更新一个Job
- GetJobList() (JobList, error) // 列出所有Job
- GetJob(name string) (Job, error) // 获取指定Job

#### Job处理
- RunJob(ctx context, name string) error // 运行Job
- StopJob(ctx context, name string) error // 停止Job

#### 可观测性
- RunningJob() ([]*JobInstance, error)
- JobHistory() ([]*JobInstance, error)

# Api
- CreateJobList(concurrency, maxScheduledCount，maxHistory int) (*JobList, error) // 创建job list同时创建事件队列
- CreateJob(name, display, type string, scheduleDuration, timeout, retryWait time.Duration, retryTimes int, jober Jober) (*Job, error)
- AddJobList(*Job) error // 添加一个Job
- DeleteJobList(name string) error // 删除一个Job
- UpdateJobList(job Job) (*Job, error) // 更新一个Job
- GetJobList() (JobList, error) // 列出所有Job
- GetJob(name string) (Job, error) // 获取指定Job
- RunJob(ctx context, name string) error // 运行Job
- StopJob(ctx context, name string) error // 停止Job
- RunningJob() ([]*JobInstance, error)
- JobHistory() ([]*JobInstance, error)

# 一个job的实例
## ShellScript
```
type ShellScript struct {
    Name string
	Version string
    ScriptTxt string
	pid int
}

func (ss *ShellScript)Run(ctx context) error {
    pid,cmd, err := runCreate(ss.ScriptTxt)
	ss.pid = pid
	err = cmd.Run()
	return err
}

func (ss *ShellScript)Run(ctx context) error {
    return os.Kill(ss.pid)
}
```

    
