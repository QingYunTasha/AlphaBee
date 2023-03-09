package main

import (
	infradomain "AlphaBee/domain/infra"
	delivery "AlphaBee/internal/delivery"
	infra "AlphaBee/internal/infra"
	usecase "AlphaBee/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	repo := infradomain.Repository{
		JobQueue:     make(chan infradomain.Job, viper.GetInt("job.queuesize")),
		TaskQueues:   make(map[string]infradomain.TaskQueue),
		Brokers:      make(map[string]infradomain.Broker),
		WorkerQueues: make(map[string]infradomain.WorkerQueue),
	}

	dispatcher := infra.NewDispatcher(repo)
	dispatcher.Run()

	alphaBeeUsecase := usecase.NewAlphaBeeUsecase(repo)

	server := gin.Default()

	delivery.NewAlphaBeeHandler(server, alphaBeeUsecase)

	server.Run(":9000")
}
