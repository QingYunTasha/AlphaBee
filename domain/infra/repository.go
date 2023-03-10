package domain

type Repository struct {
	JobQueue           chan Job
	TaskQueues         map[string]AsyncTaskQueue
	Brokers            map[string]Broker
	TaskWorkersMapping map[string]map[string]bool
	WorkerTasksMapping map[string]map[string]bool
	WorkerQueues       map[string]WorkerQueue
}
