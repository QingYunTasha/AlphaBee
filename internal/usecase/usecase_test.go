package usecase_test

import (
	infradomain "AlphaBee/domain/infra"
	infra "AlphaBee/internal/infra"
	taskqueue "AlphaBee/internal/infra/task_queue"
	usecase "AlphaBee/internal/usecase"
	"testing"
)

func TestAlphaBeeUsecase_PushJob(t *testing.T) {
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	job := infradomain.Job{
		TaskName: "job1",
		Priority: 1,
	}

	if err := usecase.PushJob(job); err != nil {
		t.Errorf("Error pushing job: %s", err.Error())
	}

	if len(repo.JobQueue) != 1 {
		t.Errorf("Expected 1 job in job queue, but found %d", len(repo.JobQueue))
	}
}

func TestAlphaBeeUsecase_PullJob(t *testing.T) {
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	repo.WorkerQueues[infradomain.WorkerName(workerName)] <- infradomain.Job{}

	if _, err := usecase.PullJob(workerName); err != nil {
		t.Errorf("Error pulling job: %s", err.Error())
	}

	if _, err := usecase.PullJob("invalid-worker"); err != nil {
		t.Errorf("expected error, but got nil")
	}
}

func TestAlphaBeeUsecase_AddTask(t *testing.T) {
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	taskName := "task1"
	algorithm := "PRIORITY_SMALL_FIRST"
	queueLength := 3
	if err := usecase.AddTask(taskName, algorithm, queueLength); err != nil {
		t.Errorf("Error add task: %s", err.Error())
	}

	if len(repo.TaskQueues) != 1 {
		t.Errorf("taskQueue length should be 1")
	}
	queues, ok := repo.TaskQueues[infradomain.TaskName(taskName)]
	if !ok {
		t.Errorf("%s queue not found", taskName)
	}
	if queues.Len() != queueLength {
		t.Errorf("len of queues should be %d", queueLength)
	}
	if _, ok := repo.Brokers[infradomain.TaskName(taskName)]; !ok {
		t.Errorf("%s broker not found", taskName)
	}

	if err := usecase.AddTask(taskName, algorithm, queueLength); err == nil {
		t.Errorf("expected error")
	}

	unValidAlgorithm := "unvalid"
	if err := usecase.AddTask(taskName, unValidAlgorithm, 1); err == nil {
		t.Errorf("expected error")
	}
}

func TestAlphaBeeUsecase_RemoveTask(t *testing.T) {
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	taskName := "task1"
	asyncTaskQueue := taskqueue.NewTaskQueue(infradomain.PrioritySmallFirst, 3)
	repo.TaskQueues[infradomain.TaskName(taskName)] = asyncTaskQueue

	broker := infra.NewBroker(asyncTaskQueue, repo.WorkerQueues)
	repo.Brokers[infradomain.TaskName(taskName)] = broker

	if err := usecase.RemoveTask(taskName); err != nil {
		t.Errorf("remove task error")
	}

	InValidTaskName := "invalid_task1"
	if err := usecase.RemoveTask(InValidTaskName); err == nil {
		t.Errorf("expect error")
	}
}

func TestAlphaBeeUsecase_AddWorkedr(t *testing.T) {
	repo := infra.NewRepository()
	usecase := usecase.NewAlphaBeeUsecase(repo)

	workerName := "worker1"
	tasks := []string{"task1", "task2", "task3"}
	queueLength := 3

	if err := usecase.AddWorker(workerName, tasks, queueLength); err != nil {
		t.Errorf("Add worker error: %v", err)
	}

	if len(repo.WorkerTasksMapping) != 1 {
		t.Errorf("lenhth of workerTasksMapping error")
	}

	tasksMap, ok := repo.WorkerTasksMapping[infradomain.WorkerName(workerName)]
	if !ok {
		t.Errorf("WorkerTasksMapping error, not found %s", workerName)
	}

	if len(tasksMap) != queueLength {
		t.Errorf("tasks length error")
	}

	for _, task := range tasks {
		if _, ok := repo.TaskWorkersMapping[infradomain.TaskName(task)]; !ok {
			t.Errorf("%s not found in workerTasksMapping", task)
		}
	}

	workerQueue, ok := repo.WorkerQueues[infradomain.WorkerName(workerName)]
	if !ok {
		t.Errorf("%s workerQueue not found")
	}
	if len(workerQueue) != queueLength {
		t.Errorf("workerQueue length error")
	}

	// TODO: check worker queue will be initialized to fill jobs

	if err := usecase.AddWorker(workerName, tasks, queueLength); err == nil {
		t.Errorf("expect error")
	}
}

func TestAlphaBeeUsecase_RemoveWorker(t *testing.T) {
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

	if err := usecase.RemoveWorker("invalid worker"); err == nil {
		t.Errorf("expect error")
	}

	if err := usecase.RemoveWorker(workerName); err != nil {
		t.Errorf("removeWorker error")
	}

	if len(repo.WorkerQueues) != 0 {
		t.Errorf("workerQueue must be zero")
	}

	for _, task := range tasks {
		if _, ok := repo.TaskWorkersMapping[infradomain.TaskName(task)]; ok {

			// not implemented
			t.Errorf("")

		}
	}
}
