package infra

import (
	infradomain "AlphaBee/domain/infra"
)

func NewWorkerQueue(id uint, workerName string, n int) *infradomain.WorkerQueue {
	return &infradomain.WorkerQueue{
		WorkerName: workerName,
		Jobs:       make(chan infradomain.Job, n),
	}
}
