package infra_test

import (
	domain "AlphaBee/domain"
	"testing"
	"time"
)

func TestDispatcher_Run(t *testing.T) {
	// Set up the dispatcher with a job queue and task queues
	jobQueue := make(chan domain.Job, 10)
	taskQueues := make(map[domain.TaskName]*domain.Queue)
	taskQueues["TaskA"] = domain.NewQueue()
	taskQueues["TaskB"] = domain.NewQueue()

	dispatcher := Dispatcher{
		JobQueue:   jobQueue,
		TaskQueues: taskQueues,
	}

	// Run the dispatcher in a separate goroutine
	go dispatcher.Run()

	// Send a job to the dispatcher's job queue
	job := domain.Job{
		TaskName: "TaskA",
		Payload:  "hello",
	}
	jobQueue <- job

	// Wait for the job to be processed
	time.Sleep(100 * time.Millisecond)

	// Check that the job was added to the correct task queue
	taskQueue := taskQueues[job.TaskName]
	if taskQueue.Size() != 1 {
		t.Errorf("Expected task queue size to be 1, but got %d", taskQueue.Size())
	}

	// Check that the job was added correctly to the task queue
	retrievedJob := taskQueue.Pop()
	if retrievedJob.TaskName != job.TaskName || retrievedJob.Payload != job.Payload {
		t.Errorf("Expected retrieved job to be %+v, but got %+v", job, retrievedJob)
	}
}
