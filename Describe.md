# Design a Job Scheduler
The goal of this exercise is to design a job scheduler that we can use to dynamically schedule different jobs. For this exercise, the scheduler needs to:
- Create/Read/Update/Delete jobs to run
- Supports one time execution and repetitive executions triggered at a fixed interval,
i.e. 10 seconds.
- Jobs are persisted
- The system is scalable to thousands of jobs and many workers
Some extra things to consider:
- Monitoring such as timeout detection?
- Retry, maybe?
- If you are using a DB, how can we make the interaction faster?
For this exercise, you can choose to use a mocked implementation for Database. For scalability, you can provide justification why your design would be feasible.
You may use Rust, Golang, Java, Nodejs for your implementation. Rust is prefered, but the rest of the languages are welcomed as well. For the exercise, please dockerize everything.
For your submission, please create a github repo and invite Konomi developers to your repo. Please ensure there is enough documentation and comments for easy readability.
