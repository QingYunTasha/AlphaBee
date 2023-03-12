package infra

import (
	infradomain "AlphaBee/domain/infra"

	"github.com/spf13/viper"
)

func NewRepository() infradomain.Repository {
	return infradomain.Repository{
		JobQueue:           make(chan infradomain.Job, viper.GetInt("job.queuesize")),
		TaskQueues:         make(map[infradomain.TaskName]infradomain.AsyncTaskQueue),
		Brokers:            make(map[infradomain.TaskName]infradomain.Broker),
		WorkerQueues:       make(map[infradomain.WorkerName]infradomain.WorkerQueue),
		TaskWorkersMapping: make(map[infradomain.TaskName]map[infradomain.WorkerName]bool),
		WorkerTasksMapping: make(map[infradomain.WorkerName]map[infradomain.TaskName]bool),
	}
}
