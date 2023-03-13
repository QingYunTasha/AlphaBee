package infra

import (
	domain "AlphaBee/domain"

	"github.com/spf13/viper"
)

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
