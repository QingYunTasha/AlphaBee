package taskqueue

import (
	infradomain "AlphaBee/domain/infra"
	"sync"
)

type TaskQueue struct {
	sync.Mutex
	Jobs infradomain.TaskQueue
}

func NewTaskQueue(algorithm infradomain.Algorithm, length int) infradomain.AsyncTaskQueue {
	var jobs infradomain.TaskQueue
	switch algorithm {
	case infradomain.PrioritySmallFirst:
		jobs = NewMinPriorityQueue()
	case infradomain.PriorityLargeFirst:
		jobs = NewMaxPriorityQueue()
	}

	return &TaskQueue{
		Mutex: sync.Mutex{},
		Jobs:  jobs,
	}
}

func (q *TaskQueue) Push(job infradomain.Job) {
	q.Lock()
	defer q.Unlock()
	q.Jobs.Push(job)
}

func (q *TaskQueue) Pop() (job infradomain.Job) {
	q.Lock()
	defer q.Unlock()
	if q.Jobs.Len() == 0 {
		return
	}

	return q.Jobs.Pop()
}

func (q *TaskQueue) Len() int {
	return q.Jobs.Len()
}
