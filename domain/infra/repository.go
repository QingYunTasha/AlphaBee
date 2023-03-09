package domain

type Repository struct {
	JobQueue     chan Job
	TaskQueues   map[string]TaskQueue
	WorkerQueues map[string]WorkerQueue
	Brokers      map[string]Broker
}
