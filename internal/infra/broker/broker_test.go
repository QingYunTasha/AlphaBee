package broker_test

import (
	"AlphaBee/domain"
	"AlphaBee/internal/infra"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroker_PushJob(t *testing.T) {
	assert := assert.New(t)
	taskQueue := infra.NewTaskQueue(domain.PriorityLargeFirst, 3)
	taskQueue.Push(domain.Job{
		TaskName: "task1",
		Priority: 1,
	})

	workerName := "worker1"
	workerQueues := make(map[domain.WorkerName]domain.WorkerQueue, 3)
	workerQueues[domain.WorkerName(workerName)] = infra.NewWorkerQueue(3)

	broker := infra.NewBroker(taskQueue, workerQueues)

	broker.PushJob(workerName)

	assert.Equal(1, len(workerQueues[domain.WorkerName(workerName)]))
}
