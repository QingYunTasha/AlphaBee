package main

import (
	"AlphaBee/config"
	infradomain "AlphaBee/domain/infra"
	delivery "AlphaBee/internal/delivery"
	infra "AlphaBee/internal/infra"
	usecase "AlphaBee/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	repo := infradomain.Repository{
		JobQueue:           make(chan infradomain.Job, viper.GetInt("job.queuesize")),
		TaskQueues:         make(map[string]infradomain.AsyncTaskQueue),
		Brokers:            make(map[string]infradomain.Broker),
		WorkerQueues:       make(map[string]infradomain.WorkerQueue),
		TaskWorkersMapping: make(map[string]map[string]infradomain.WorkerQueue),
	}

	dispatcher := infra.NewDispatcher(repo.JobQueue, repo.TaskQueues)
	dispatcher.Run()

	alphaBeeUsecase := usecase.NewAlphaBeeUsecase(repo)

	server := gin.Default()

	delivery.NewAlphaBeeHandler(server, alphaBeeUsecase)

	server.Run(":9000")
}
