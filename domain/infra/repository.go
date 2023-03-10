package domain

type Repository struct {
	JobQueue           chan Job
	TaskQueues         map[string]AsyncTaskQueue
	Brokers            map[string]Broker
	TaskWorkersMapping map[string]map[string]WorkerQueue
	WorkerQueues       map[string]WorkerQueue
}
