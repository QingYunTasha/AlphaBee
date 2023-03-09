package infra

import (
	infradomain "AlphaBee/domain/infra"
)

type Broker struct {
	TaskQueue infradomain.TaskQueue
	Repo      infradomain.Repository
}

func NewBroker(taskQueue infradomain.TaskQueue) infradomain.Broker {
	return &Broker{
		TaskQueue: taskQueue,
	}
}

func (b Broker) PushJob(job infradomain.Job, workerName string) {
	j := b.TaskQueue.Pop()
	b.Repo.WorkerQueues[workerName].Jobs <- j
}
