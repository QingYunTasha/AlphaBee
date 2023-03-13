package taskqueue

import (
	infradomain "AlphaBee/domain/infra"
)

type MaxPriorityQueue []infradomain.Job

func NewMaxPriorityQueue() infradomain.TaskQueue {
	q := append(MaxPriorityQueue{}, infradomain.Job{})
	return &q
}

func (q MaxPriorityQueue) Len() int { return len(q) - 1 }

func (q MaxPriorityQueue) less(i, j int) bool { return q[i].Priority < q[j].Priority }

func (q MaxPriorityQueue) swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q *MaxPriorityQueue) Push(job infradomain.Job) {
	*q = append(*q, job)
	q.up(len(*q) - 1)
}

func (q *MaxPriorityQueue) Pop() infradomain.Job {
	n := len(*q) - 1
	res := (*q)[1]
	(*q)[1] = (*q)[n]
	*q = (*q)[:n]
	q.down(1)
	return res
}

func (q MaxPriorityQueue) up(i int) {
	for j := i >> 1; i > 1 && q.less(j, i); i, j = j, j>>1 {
		q.swap(i, j)
	}
}

func (q MaxPriorityQueue) down(i int) {
	for j := i << 1; j < len(q); i, j = j, j<<1 {
		if j+1 < len(q) && q.less(j, j+1) {
			j++
		}
		if q.less(j, i) {
			break
		}
		q.swap(i, j)
	}
}
