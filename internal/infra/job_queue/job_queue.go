package infra

import (
	domain "AlphaBee/domain"
)

func NewJobQueue(n int) domain.JobQueue {
	return make(domain.JobQueue, n)
}
