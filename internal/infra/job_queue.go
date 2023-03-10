package infra

import (
	infradomain "AlphaBee/domain/infra"
)

func NewJobQueue(n int) infradomain.JobQueue {
	return make(infradomain.JobQueue, n)
}
