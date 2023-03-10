package delivery

import (
	infradomain "AlphaBee/domain/infra"
	usecasedomain "AlphaBee/domain/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AlphaBeeHandler struct {
	AlphaBeeUsecase usecasedomain.AlphaBeeUsecase
}

func NewAlphaBeeHandler(server *gin.Engine, usecase usecasedomain.AlphaBeeUsecase) {
	handler := &AlphaBeeHandler{
		AlphaBeeUsecase: usecase,
	}

	server.POST("/pushjob", handler.PushJob)
	server.POST("/pulljob/:workername", handler.PopJob)
	server.POST("/task", handler.AddTask)
	server.DELETE("/task/:name", handler.RemoveTask)
	server.POST("/worker", handler.AddWorker)
	server.DELETE("/worker/:name", handler.RemoveWorker)
}

func (h AlphaBeeHandler) PushJob(c *gin.Context) {
	var job infradomain.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AlphaBeeUsecase.PushJob(job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h AlphaBeeHandler) PopJob(c *gin.Context) {
	job, err := h.AlphaBeeUsecase.PopJob(c.Param("workername"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h AlphaBeeHandler) AddTask(c *gin.Context) {
	var body struct {
		TaskName  string                `json:"task_name" binding:"required"`
		Algorithm infradomain.Algorithm `json:"algorithm" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AlphaBeeUsecase.AddTask(body.TaskName, body.Algorithm, viper.GetInt("task.queuesize")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h AlphaBeeHandler) RemoveTask(c *gin.Context) {
	if err := h.AlphaBeeUsecase.RemoveTask(c.Param("name")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h AlphaBeeHandler) AddWorker(c *gin.Context) {
	var body struct {
		WorkerName string   `json:"worker_name" binding:"required"`
		Tasks      []string `json:"tasks" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AlphaBeeUsecase.AddWorker(body.WorkerName, body.Tasks, viper.GetInt("worker.queuesize")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h AlphaBeeHandler) RemoveWorker(c *gin.Context) {
	if err := h.AlphaBeeUsecase.RemoveWorker(c.Param("name")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
