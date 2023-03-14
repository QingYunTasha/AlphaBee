package dispatcher_test

import (
	"AlphaBee/domain"
	"AlphaBee/internal/infra"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDispatcher_Run(t *testing.T) {
	assert := assert.New(t)
	// Set up the dispatcher with a job queue and task queues
	jobQueue := make(chan domain.Job, 10)
	taskQueues := make(map[domain.TaskName]domain.AsyncTaskQueue)
	taskQueues["TaskA"] = infra.NewTaskQueue(domain.PriorityLargeFirst, 3)
	taskQueues["TaskB"] = infra.NewTaskQueue(domain.PriorityLargeFirst, 3)

	dispatcher := infra.NewDispatcher(jobQueue, taskQueues)

	// Run the dispatcher in a separate goroutine
	go dispatcher.Run()

	// Send a job to the dispatcher's job queue
	job := domain.Job{
		TaskName: "TaskA",
		Priority: 1,
	}
	jobQueue <- job

	job = domain.Job{
		TaskName: "TaskB",
		Priority: 1,
	}
	jobQueue <- job

	// Wait for the job to be processed
	time.Sleep(100 * time.Millisecond)

	assert.Equal(1, taskQueues["TaskA"].Len())
	assert.Equal(1, taskQueues["TaskB"].Len())

}
