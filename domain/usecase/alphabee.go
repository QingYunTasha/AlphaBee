package domain

import (
	dbdomain "AlphaBee/domain/infra/database"
)

type AlphabeeUsecase interface {
	Push(dbdomain.Job) error
	Pop(workerID uint) (dbdomain.Job, error)
	AddTask(dbdomain.Task) error
	RemoveTask(taskName string) error
	AddWorker(dbdomain.Worker) error
	RemoveWorker(workerName string) error
}

type TaskDispatcher interface {
	AssignJobToTask()
}

type WorkerDispatcher interface {
	AssignJobToWorker()
}
