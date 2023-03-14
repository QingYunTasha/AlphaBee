package dispatcher

import (
	domain "AlphaBee/domain"
)

type Dispatcher struct {
	JobQueue   domain.JobQueue
	TaskQueues map[domain.TaskName]domain.AsyncTaskQueue
}

func NewDispatcher(jobQueue domain.JobQueue, taskQueues map[domain.TaskName]domain.AsyncTaskQueue) domain.Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		TaskQueues: taskQueues,
	}
}

func (d Dispatcher) Run() {
	go func() {
		var job domain.Job
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
