package infra

import (
	infradomain "AlphaBee/domain/infra"
)

type Dispatcher struct {
	JobQueue   infradomain.JobQueue
	TaskQueues map[infradomain.TaskName]infradomain.AsyncTaskQueue
}

func NewDispatcher(jobQueue infradomain.JobQueue, taskQueues map[infradomain.TaskName]infradomain.AsyncTaskQueue) infradomain.Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		TaskQueues: taskQueues,
	}
}

func (d Dispatcher) Run() {
	go func() {
		var job infradomain.Job
		for {
			job = <-d.JobQueue

			for key := range d.TaskQueues {
				if job.TaskName == string(key) {
					d.TaskQueues[key].Push(job)
					break
				}
			}
		}
	}()
}
