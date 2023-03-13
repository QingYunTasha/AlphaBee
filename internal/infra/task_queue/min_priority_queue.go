package taskqueue

import (
	domain "AlphaBee/domain"
)

type MinPriorityQueue []domain.Job

func NewMinPriorityQueue() domain.TaskQueue {
	q := append(MinPriorityQueue{}, domain.Job{})
	return &q
}

func (q MinPriorityQueue) Len() int { return len(q) - 1 }

func (q MinPriorityQueue) less(i, j int) bool { return q[i].Priority > q[j].Priority }

func (q MinPriorityQueue) swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q *MinPriorityQueue) Push(job domain.Job) {
	*q = append(*q, job)
	q.up(len(*q) - 1)
}

func (q *MinPriorityQueue) Pop() domain.Job {
	n := len(*q) - 1
	res := (*q)[1]
	(*q)[1] = (*q)[n]
	*q = (*q)[:n]
	q.down(1)
	return res
}

func (q MinPriorityQueue) up(i int) {
	for j := i >> 1; i > 1 && q.less(j, i); i, j = j, j>>1 {
		q.swap(i, j)
	}
}

func (q MinPriorityQueue) down(i int) {
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
