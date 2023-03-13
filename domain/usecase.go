package domain

type AlphaBeeUsecase interface {
	PushJob(Job) error
	PullJob(workerName string) (Job, error)
	AddTask(taskName string, algorithm string, taskQueueLength int) error
	RemoveTask(taskName string) error
	AddWorker(workerName string, tasks []string, workerQueueLength int) error
	RemoveWorker(workerName string) error
}
