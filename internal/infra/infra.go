package infra

import (
	domain "AlphaBee/domain"
	"AlphaBee/internal/infra/broker"
	"AlphaBee/internal/infra/dispatcher"
	taskqueue "AlphaBee/internal/infra/task_queue"
	"sync"

	"github.com/spf13/viper"
)

type Infra struct{}

func (i Infra) NewBroker(taskQueue domain.AsyncTaskQueue, workerQueues map[domain.WorkerName]domain.WorkerQueue) domain.Broker {
	return &broker.Broker{
		TaskQueue:    taskQueue,
		WorkerQueues: workerQueues,
	}
}

func NewWorkerQueue(n int) domain.WorkerQueue {
	return make(domain.WorkerQueue, n)
}

func NewTaskQueue(algorithm domain.Algorithm, length int) domain.AsyncTaskQueue {
	var jobs domain.TaskQueue
	switch algorithm {
	case domain.PrioritySmallFirst:
		jobs = taskqueue.NewMinPriorityQueue()
	case domain.PriorityLargeFirst:
		jobs = taskqueue.NewMaxPriorityQueue()
	}

	return &taskqueue.TaskQueue{
		Mutex: sync.Mutex{},
		Jobs:  jobs,
	}
}

func NewRepository() *domain.Repository {
	return &domain.Repository{
		JobQueue:           make(chan domain.Job, viper.GetInt("job.queuesize")),
		TaskQueues:         make(map[domain.TaskName]domain.AsyncTaskQueue),
		Brokers:            make(map[domain.TaskName]domain.Broker),
		WorkerQueues:       make(map[domain.WorkerName]domain.WorkerQueue),
		TaskWorkersMapping: make(map[domain.TaskName]map[domain.WorkerName]bool),
		WorkerTasksMapping: make(map[domain.WorkerName]map[domain.TaskName]bool),
	}
}

func NewJobQueue(n int) domain.JobQueue {
	return make(domain.JobQueue, n)
}

func NewDispatcher(jobQueue domain.JobQueue, taskQueues map[domain.TaskName]domain.AsyncTaskQueue) domain.Dispatcher {
	return &dispatcher.Dispatcher{
		JobQueue:   jobQueue,
		TaskQueues: taskQueues,
	}
}

func NewBroker(taskQueue domain.AsyncTaskQueue, workerQueues map[domain.WorkerName]domain.WorkerQueue) domain.Broker {
	return &broker.Broker{
		TaskQueue:    taskQueue,
		WorkerQueues: workerQueues,
	}
}
