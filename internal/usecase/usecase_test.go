package usecase_test

import (
	domain "AlphaBee/domain"
	infra "AlphaBee/internal/infra"
	usecase "AlphaBee/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: must structured, ref: https://ithelp.ithome.com.tw/articles/10204692

func TestAlphaBeeUsecase_PushJob(t *testing.T) {
	assert := assert.New(t)
	repo := infra.NewRepository()
	repo.JobQueue = make(chan domain.Job, 1)
	usecase := usecase.NewAlphaBeeUsecase(repo)

	job := domain.Job{
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
	repo.JobQueue = make(chan domain.Job, 1)
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	repo.WorkerQueues[domain.WorkerName(workerName)] = make(domain.WorkerQueue, 1)
	repo.WorkerQueues[domain.WorkerName(workerName)] <- domain.Job{}

	_, err := usecase.PullJob(workerName)
	assert.Nil(err)
	assert.Equal(0, len(repo.WorkerQueues[domain.WorkerName(workerName)]))

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

	_, ok := repo.TaskQueues[domain.TaskName(taskName)]
	assert.True(ok)
	// priority queue is not fixed length
	//assert.Equal(queueLength, queues.Len())

	_, ok = repo.Brokers[domain.TaskName(taskName)]
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
	asyncTaskQueue := infra.NewTaskQueue(domain.PrioritySmallFirst, 3)
	repo.TaskQueues[domain.TaskName(taskName)] = asyncTaskQueue

	broker := infra.NewBroker(asyncTaskQueue, repo.WorkerQueues)
	repo.Brokers[domain.TaskName(taskName)] = broker

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
		repo.TaskQueues[domain.TaskName(task)] = infra.NewTaskQueue(domain.PrioritySmallFirst, 3)
	}
	queueLength := 3

	err := usecase.AddWorker(workerName, tasks, queueLength)
	assert.Nil(err)
	assert.Equal(1, len(repo.WorkerTasksMapping))

	tasksMap, ok := repo.WorkerTasksMapping[domain.WorkerName(workerName)]
	assert.True(ok)
	assert.Equal(queueLength, len(tasksMap))

	for _, task := range tasks {
		_, ok := repo.TaskWorkersMapping[domain.TaskName(task)]
		assert.True(ok)
	}

	workerQueue, ok := repo.WorkerQueues[domain.WorkerName(workerName)]
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
	repo.WorkerQueues[domain.WorkerName(workerName)] = wq

	for _, task := range tasks {
		repo.TaskWorkersMapping[domain.TaskName(task)] = make(map[domain.WorkerName]bool)
	}

	err := usecase.RemoveWorker("invalid worker")
	assert.NotNil(err)

	err = usecase.RemoveWorker(workerName)
	assert.Nil(err)
	assert.Equal(0, len(repo.WorkerQueues))

	for _, task := range tasks {
		workers, ok := repo.TaskWorkersMapping[domain.TaskName(task)]
		assert.True(ok)

		_, ok = workers[domain.WorkerName(workerName)]
		assert.False(ok)
	}
}
