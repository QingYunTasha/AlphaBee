package infra

import (
	infradomain "AlphaBee/domain/infra"
)

type Broker struct {
	TaskQueue infradomain.AsyncTaskQueue
	Repo      infradomain.Repository
}

func NewBroker(taskQueue infradomain.AsyncTaskQueue) infradomain.Broker {
	return &Broker{
		TaskQueue: taskQueue,
	}
}

func (b Broker) PushJob(job infradomain.Job, workerName string) {
	j := b.TaskQueue.Pop()
	b.Repo.WorkerQueues[workerName] <- j
}
