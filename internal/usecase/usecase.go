package usecase

import (
	infradomain "AlphaBee/domain/infra"
	usecasedomain "AlphaBee/domain/usecase"
	"AlphaBee/internal/infra"
	taskqueue "AlphaBee/internal/infra/task_queue"
	"fmt"
)

type AlphaBeeUsecase struct {
	repo infradomain.Repository
}

func NewAlphaBeeUsecase(repo infradomain.Repository) usecasedomain.AlphaBeeUsecase {
	return &AlphaBeeUsecase{
		repo: repo,
	}
}

func (a AlphaBeeUsecase) PushJob(job infradomain.Job) error {
	a.repo.JobQueue <- job
	return nil
}

func (a AlphaBeeUsecase) PullJob(workerName string) (infradomain.Job, error) {
	wq, ok := a.repo.WorkerQueues[infradomain.WorkerName(workerName)]
	if !ok {
		return infradomain.Job{}, fmt.Errorf("worker %s not found", workerName)
	}

	job := <-wq

	go func() {
	LOOP:
		for task := range a.repo.WorkerTasksMapping[infradomain.WorkerName(workerName)] {
			if a.repo.TaskQueues[task].Len() > 0 {
				wq <- a.repo.TaskQueues[task].Pop()
				break LOOP
			}
		}
	}()

	return job, nil
}

func (a AlphaBeeUsecase) AddTask(taskName string, algorithm string, n int) error {
	if _, ok := a.repo.TaskQueues[infradomain.TaskName(taskName)]; ok {
		return fmt.Errorf("task %s already exists", taskName)
	}

	if isValid := infradomain.IsValidAlgorithm(algorithm); !isValid {
		return fmt.Errorf("algorithm %s not supported", algorithm)
	}

	tq := taskqueue.NewTaskQueue(infradomain.Algorithm(algorithm), n)
	a.repo.TaskQueues[infradomain.TaskName(taskName)] = tq
	a.repo.Brokers[infradomain.TaskName(taskName)] = infra.NewBroker(tq, a.repo.WorkerQueues)

	return nil
}

func (a AlphaBeeUsecase) RemoveTask(taskName string) error {
	if _, ok := a.repo.TaskQueues[infradomain.TaskName(taskName)]; !ok {
		return fmt.Errorf("task %s not found", taskName)
	}
	delete(a.repo.TaskQueues, infradomain.TaskName(taskName))
	delete(a.repo.Brokers, infradomain.TaskName(taskName))
	return nil
}

func (a AlphaBeeUsecase) AddWorker(workerName string, taskNames []string, n int) error {
	if _, ok := a.repo.WorkerQueues[infradomain.WorkerName(workerName)]; ok {
		return fmt.Errorf("worker %s already exists", workerName)
	}

	// add m2m mapping
	if _, ok := a.repo.WorkerTasksMapping[infradomain.WorkerName(workerName)]; !ok {
		a.repo.WorkerTasksMapping[infradomain.WorkerName(workerName)] = make(map[infradomain.TaskName]bool)
	}

	for _, taskName := range taskNames {
		if _, ok := a.repo.TaskWorkersMapping[infradomain.TaskName(taskName)]; !ok {
			a.repo.TaskWorkersMapping[infradomain.TaskName(taskName)] = make(map[infradomain.WorkerName]bool)
		}
		a.repo.TaskWorkersMapping[infradomain.TaskName(taskName)][infradomain.WorkerName(workerName)] = true
		a.repo.WorkerTasksMapping[infradomain.WorkerName(workerName)][infradomain.TaskName(taskName)] = true
	}

	wq := infra.NewWorkerQueue(n)
	go func() {
	LOOP:
		for task := range a.repo.WorkerTasksMapping[infradomain.WorkerName(workerName)] {
			for a.repo.TaskQueues[task].Len() > 0 {
				wq <- a.repo.TaskQueues[task].Pop()

				if len(wq) == cap(wq) {
					break LOOP
				}
			}
		}
	}()
	a.repo.WorkerQueues[infradomain.WorkerName(workerName)] = wq
	return nil
}

func (a AlphaBeeUsecase) RemoveWorker(workerName string) error {
	if _, ok := a.repo.WorkerQueues[infradomain.WorkerName(workerName)]; !ok {
		return fmt.Errorf("worker %s not found", workerName)
	}

	// TODO: This method is very inefficient, try to find another way to
	// store worker - task mappings
	for _, workers := range a.repo.TaskWorkersMapping {
		delete(workers, infradomain.WorkerName(workerName))
	}

	delete(a.repo.WorkerQueues, infradomain.WorkerName(workerName))
	return nil
}
