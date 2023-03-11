package domain

type Repository struct {
	JobQueue           chan Job
	TaskQueues         map[TaskName]AsyncTaskQueue
	Brokers            map[TaskName]Broker
	TaskWorkersMapping map[TaskName]map[WorkerName]bool
	WorkerTasksMapping map[WorkerName]map[TaskName]bool
	WorkerQueues       map[WorkerName]WorkerQueue
}
