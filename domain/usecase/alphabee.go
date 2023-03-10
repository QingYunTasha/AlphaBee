package domain

import (
	infradomain "AlphaBee/domain/infra"
)

type AlphaBeeUsecase interface {
	PushJob(infradomain.Job) error
	PopJob(workerName string) (infradomain.Job, error)
	AddTask(taskName string, algorithm infradomain.Algorithm, n int) error
	RemoveTask(taskName string) error
	AddWorker(workerName string, tasks []string, n int) error
	RemoveWorker(workerName string) error
}
