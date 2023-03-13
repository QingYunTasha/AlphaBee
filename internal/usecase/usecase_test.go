package usecase_test

import (
	infradomain "AlphaBee/domain/infra"
	infra "AlphaBee/internal/infra"
	taskqueue "AlphaBee/internal/infra/task_queue"
	usecase "AlphaBee/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: must structured, ref: https://ithelp.ithome.com.tw/articles/10204692

func TestAlphaBeeUsecase_PushJob(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	repo.JobQueue = make(chan infradomain.Job, 1)
	usecase := usecase.NewAlphaBeeUsecase(repo)

	job := infradomain.Job{
		TaskName: "job1",
		Priority: 1,
	}

	err := usecase.PushJob(job)
	assert.Nil(err)

	assert.Equal(1, len(repo.JobQueue))
}

func TestAlphaBeeUsecase_PullJob(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	repo.JobQueue = make(chan infradomain.Job, 1)
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	repo.WorkerQueues[infradomain.WorkerName(workerName)] = make(infradomain.WorkerQueue, 1)
	repo.WorkerQueues[infradomain.WorkerName(workerName)] <- infradomain.Job{}

	_, err := usecase.PullJob(workerName)
	assert.Nil(err)
	assert.Equal(0, len(repo.WorkerQueues[infradomain.WorkerName(workerName)]))

	_, err = usecase.PullJob("invalid-worker")
	assert.NotNil(err)
}

func TestAlphaBeeUsecase_AddTask(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	err := usecase.AddTask("", "invalid-algorithm", 1)
	assert.NotNil(err)

	taskName := "task1"
	algorithm := "PRIORITY_SMALL_FIRST"
	queueLength := 3

	err = usecase.AddTask(taskName, algorithm, queueLength)
	assert.Nil(err)
	assert.Equal(1, len(repo.TaskQueues))

	_, ok := repo.TaskQueues[infradomain.TaskName(taskName)]
	assert.True(ok)
	// priority queue is not fixed length
	//assert.Equal(queueLength, queues.Len())

	_, ok = repo.Brokers[infradomain.TaskName(taskName)]
	assert.True(ok)

	err = usecase.AddTask(taskName, algorithm, queueLength)
	assert.NotNil(err)

	unValidAlgorithm := "unvalid"
	err = usecase.AddTask(taskName, unValidAlgorithm, 1)
	assert.NotNil(err)
}

func TestAlphaBeeUsecase_RemoveTask(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	taskName := "task1"
	asyncTaskQueue := taskqueue.NewTaskQueue(infradomain.PrioritySmallFirst, 3)
	repo.TaskQueues[infradomain.TaskName(taskName)] = asyncTaskQueue

	broker := infra.NewBroker(asyncTaskQueue, repo.WorkerQueues)
	repo.Brokers[infradomain.TaskName(taskName)] = broker

	err := usecase.RemoveTask(taskName)
	assert.Nil(err)
	assert.Equal(0, len(repo.TaskQueues))
	assert.Equal(0, len(repo.Brokers))

	InValidTaskName := "invalid_task1"
	err = usecase.RemoveTask(InValidTaskName)
	assert.NotNil(err)
}

func TestAlphaBeeUsecase_AddWorkedr(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	tasks := []string{"task1", "task2", "task3"}
	for _, task := range tasks {
		repo.TaskQueues[infradomain.TaskName(task)] = taskqueue.NewTaskQueue(infradomain.PrioritySmallFirst, 3)
	}
	queueLength := 3

	err := usecase.AddWorker(workerName, tasks, queueLength)
	assert.Nil(err)
	assert.Equal(1, len(repo.WorkerTasksMapping))

	tasksMap, ok := repo.WorkerTasksMapping[infradomain.WorkerName(workerName)]
	assert.True(ok)
	assert.Equal(queueLength, len(tasksMap))

	for _, task := range tasks {
		_, ok := repo.TaskWorkersMapping[infradomain.TaskName(task)]
		assert.True(ok)
	}

	workerQueue, ok := repo.WorkerQueues[infradomain.WorkerName(workerName)]
	assert.True(ok)
	assert.Equal(queueLength, cap(workerQueue))

	// TODO: check worker queue will be initialized to fill jobs

	err = usecase.AddWorker(workerName, tasks, queueLength)
	assert.NotNil(err)
}

func TestAlphaBeeUsecase_RemoveWorker(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	tasks := []string{"task1", "task2", "task3"}
	queueLength := 3

	wq := infra.NewWorkerQueue(queueLength)
	repo.WorkerQueues[infradomain.WorkerName(workerName)] = wq

	for _, task := range tasks {
		repo.TaskWorkersMapping[infradomain.TaskName(task)] = make(map[infradomain.WorkerName]bool)
	}

	err := usecase.RemoveWorker("invalid worker")
	assert.NotNil(err)

	err = usecase.RemoveWorker(workerName)
	assert.Nil(err)
	assert.Equal(0, len(repo.WorkerQueues))

	for _, task := range tasks {
		workers, ok := repo.TaskWorkersMapping[infradomain.TaskName(task)]
		assert.True(ok)

		_, ok = workers[infradomain.WorkerName(workerName)]
		assert.False(ok)
	}
}
