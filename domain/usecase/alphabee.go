package domain

import (
	dbdomain "AlphaBee/domain/infra/database"
)

type AlphabeeUsecase interface {
	Push(dbdomain.Job) error
	Pop(workerID uint) (dbdomain.Job, error)
	AddTask(taskName string, algorithm dbdomain.Algorithm) error
	RemoveTask(taskName string) error
	AddWorker(workerName string) error
	RemoveWorker(workerName string) error
}

type DispatcherFactory interface {
	AddDispatcher() error
	RemoveDispatcher() error
}

type BrokerFactory interface {
	AddBroker() error
	RemoveBroker() error
}
