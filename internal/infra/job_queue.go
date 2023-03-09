package infra

import (
	infradomain "AlphaBee/domain/infra"
)

func NewJobQueue(n int) chan infradomain.Job {
	return make(chan infradomain.Job, n)
}
