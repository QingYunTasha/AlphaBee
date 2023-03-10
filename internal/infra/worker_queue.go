package infra

import (
	infradomain "AlphaBee/domain/infra"
)

func NewWorkerQueue(n int) infradomain.WorkerQueue {
	return make(infradomain.WorkerQueue, n)
}
