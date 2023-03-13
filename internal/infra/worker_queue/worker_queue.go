package infra

import (
	domain "AlphaBee/domain"
)

func NewWorkerQueue(n int) domain.WorkerQueue {
	return make(domain.WorkerQueue, n)
}
