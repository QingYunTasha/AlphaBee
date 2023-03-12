package main

import (
	"AlphaBee/config"
	delivery "AlphaBee/internal/delivery/restful"
	infra "AlphaBee/internal/infra"
	usecase "AlphaBee/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	repo := infra.NewRepository()

	dispatcher := infra.NewDispatcher(repo.JobQueue, repo.TaskQueues)
	dispatcher.Run()

	alphaBeeUsecase := usecase.NewAlphaBeeUsecase(repo)

	server := gin.Default()

	delivery.NewAlphaBeeHandler(server, alphaBeeUsecase)

	server.Run(":9000")
}
