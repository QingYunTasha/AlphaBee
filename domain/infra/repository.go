package domain

type Repository struct {
	JobQueue     chan Job
	TaskQueues   map[string]AsyncTaskQueue
	WorkerQueues map[string]WorkerQueue
	Brokers      map[string]Broker
}
