package domain

type AlphaBeeUsecase interface {
	PushJob(Job) error
	PopJob(workerName string) (Job, error)
	AddTask(taskName string, algorithm Algorithm, n int) error
	RemoveTask(taskName string) error
	AddWorker(workerName string, n int) error
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

type Broker interface {
	PushJob(job Job, workerName string)
}

type Dispatcher interface {
	Run()
}
