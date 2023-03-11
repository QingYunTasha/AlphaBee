package main

import (
	"AlphaBee/config"
	infradomain "AlphaBee/domain/infra"
	delivery "AlphaBee/internal/delivery/restful"
	infra "AlphaBee/internal/infra"
	usecase "AlphaBee/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	repo := infradomain.Repository{
		JobQueue:           make(chan infradomain.Job, viper.GetInt("job.queuesize")),
		TaskQueues:         make(map[infradomain.TaskName]infradomain.AsyncTaskQueue),
		Brokers:            make(map[infradomain.TaskName]infradomain.Broker),
		WorkerQueues:       make(map[infradomain.WorkerName]infradomain.WorkerQueue),
		TaskWorkersMapping: make(map[infradomain.TaskName]map[infradomain.WorkerName]bool),
		WorkerTasksMapping: make(map[infradomain.WorkerName]map[infradomain.TaskName]bool),
	}

	dispatcher := infra.NewDispatcher(repo.JobQueue, repo.TaskQueues)
	dispatcher.Run()

	alphaBeeUsecase := usecase.NewAlphaBeeUsecase(repo)

	server := gin.Default()

	delivery.NewAlphaBeeHandler(server, alphaBeeUsecase)

	server.Run(":9000")
}
