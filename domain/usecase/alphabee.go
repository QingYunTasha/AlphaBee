package domain

import (
	infradomain "AlphaBee/domain/infra"
)

type AlphaBeeUsecase interface {
	PushJob(infradomain.Job) error
	PullJob(workerName string) (infradomain.Job, error)
	AddTask(taskName string, algorithm string, taskQueueLength int) error
	RemoveTask(taskName string) error
	AddWorker(workerName string, tasks []string, workerQueueLength int) error
	RemoveWorker(workerName string) error
}
