package infra

import (
	domain "AlphaBee/domain"
)

type Broker struct {
	TaskQueue    domain.AsyncTaskQueue
	WorkerQueues map[domain.WorkerName]domain.WorkerQueue
}

func NewBroker(taskQueue domain.AsyncTaskQueue, workerQueues map[domain.WorkerName]domain.WorkerQueue) domain.Broker {
	return &Broker{
		TaskQueue:    taskQueue,
		WorkerQueues: workerQueues,
	}
}

func (b Broker) PushJob(job domain.Job, workerName string) {
	b.WorkerQueues[domain.WorkerName(workerName)] <- b.TaskQueue.Pop()
}
