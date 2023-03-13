package domain

import "sync"

type Broker interface {
	PushJob(job Job, workerName string)
}

type Dispatcher interface {
	Run()
}

type JobQueue chan Job

type Job struct {
	TaskName string
	Priority uint8
}

type PriorityJob struct {
	TaskName string
	Priority uint8
}

type Repository struct {
	JobQueue           chan Job
	TaskQueues         map[TaskName]AsyncTaskQueue
	Brokers            map[TaskName]Broker
	TaskWorkersMapping map[TaskName]map[WorkerName]bool
	WorkerTasksMapping map[WorkerName]map[TaskName]bool
	WorkerQueues       map[WorkerName]WorkerQueue
}

type Algorithm string

var (
	RoundRobin         Algorithm = "ROUND_ROBIN"
	PrioritySmallFirst Algorithm = "PRIORITY_SMALL_FIRST"
	PriorityLargeFirst Algorithm = "PRIORITY_LARGE_FIRST"
	ShortestJobFirst   Algorithm = "SHORTEST_JOB_FIRST"
)

func IsValidAlgorithm(input string) bool {
	switch Algorithm(input) {
	case RoundRobin, PriorityLargeFirst, PrioritySmallFirst, ShortestJobFirst:
		return true
	default:
		return false
	}
}

type TaskName string

type TaskQueue interface {
	Push(job Job)
	Pop() (job Job)
	Len() (n int)
}

type AsyncTaskQueue interface {
	sync.Locker
	Push(job Job)
	Pop() (job Job)
	Len() int
}

type WorkerName string

type WorkerQueue chan Job
