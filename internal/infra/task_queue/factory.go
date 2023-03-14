package taskqueue

import (
	domain "AlphaBee/domain"
	"sync"
)

type TaskQueue struct {
	sync.Mutex
	Jobs domain.TaskQueue
}

func (q *TaskQueue) Push(job domain.Job) {
	q.Lock()
	defer q.Unlock()
	q.Jobs.Push(job)
}

func (q *TaskQueue) Pop() (job domain.Job) {
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
