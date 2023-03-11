package infra

import (
	infradomain "AlphaBee/domain/infra"
)

type Broker struct {
	TaskQueue    infradomain.AsyncTaskQueue
	WorkerQueues map[infradomain.WorkerName]infradomain.WorkerQueue
}

func NewBroker(taskQueue infradomain.AsyncTaskQueue, workerQueues map[infradomain.WorkerName]infradomain.WorkerQueue) infradomain.Broker {
	return &Broker{
		TaskQueue:    taskQueue,
		WorkerQueues: workerQueues,
	}
}

func (b Broker) PushJob(job infradomain.Job, workerName string) {
	b.WorkerQueues[infradomain.WorkerName(workerName)] <- b.TaskQueue.Pop()
}
