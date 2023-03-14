package broker

import (
	domain "AlphaBee/domain"
)

type Broker struct {
	TaskQueue    domain.AsyncTaskQueue
	WorkerQueues map[domain.WorkerName]domain.WorkerQueue
}

func (b Broker) PushJob(job domain.Job, workerName string) {
	b.WorkerQueues[domain.WorkerName(workerName)] <- b.TaskQueue.Pop()
}
