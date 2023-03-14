package usecase

import (
	"AlphaBee/domain"
	"AlphaBee/internal/infra"
	"fmt"
)

type AlphaBeeUsecase struct {
	repo *domain.Repository
}

func NewAlphaBeeUsecase(repo *domain.Repository) domain.AlphaBeeUsecase {
	return &AlphaBeeUsecase{
		repo: repo,
	}
}

func (a AlphaBeeUsecase) PushJob(job domain.Job) error {
	a.repo.JobQueue <- job
	return nil
}

func (a AlphaBeeUsecase) PullJob(workerName string) (domain.Job, error) {
	wq, ok := a.repo.WorkerQueues[domain.WorkerName(workerName)]
	if !ok {
		return domain.Job{}, fmt.Errorf("worker %s not found", workerName)
	}

	job := <-wq

	go func() {
	LOOP:
		for task := range a.repo.WorkerTasksMapping[domain.WorkerName(workerName)] {
			if a.repo.TaskQueues[task].Len() > 0 {
				wq <- a.repo.TaskQueues[task].Pop()
				break LOOP
			}
		}
	}()

	return job, nil
}

func (a AlphaBeeUsecase) AddTask(taskName string, algorithm string, taskQueueLength int) error {
	if _, ok := a.repo.TaskQueues[domain.TaskName(taskName)]; ok {
		return fmt.Errorf("task %s already exists", taskName)
	}

	if isValid := domain.IsValidAlgorithm(algorithm); !isValid {
		return fmt.Errorf("algorithm %s not supported", algorithm)
	}

	tq := infra.NewTaskQueue(domain.Algorithm(algorithm), taskQueueLength)
	a.repo.TaskQueues[domain.TaskName(taskName)] = tq
	a.repo.Brokers[domain.TaskName(taskName)] = infra.NewBroker(tq, a.repo.WorkerQueues)

	return nil
}

func (a AlphaBeeUsecase) RemoveTask(taskName string) error {
	if _, ok := a.repo.TaskQueues[domain.TaskName(taskName)]; !ok {
		return fmt.Errorf("task %s not found", taskName)
	}
	delete(a.repo.TaskQueues, domain.TaskName(taskName))
	delete(a.repo.Brokers, domain.TaskName(taskName))
	return nil
}

func (a AlphaBeeUsecase) AddWorker(workerName string, taskNames []string, workerQueueLength int) error {
	if _, ok := a.repo.WorkerQueues[domain.WorkerName(workerName)]; ok {
		return fmt.Errorf("worker %s already exists", workerName)
	}

	// add m2m mapping
	if _, ok := a.repo.WorkerTasksMapping[domain.WorkerName(workerName)]; !ok {
		a.repo.WorkerTasksMapping[domain.WorkerName(workerName)] = make(map[domain.TaskName]bool)
	}

	for _, taskName := range taskNames {
		if _, ok := a.repo.TaskWorkersMapping[domain.TaskName(taskName)]; !ok {
			a.repo.TaskWorkersMapping[domain.TaskName(taskName)] = make(map[domain.WorkerName]bool)
		}
		a.repo.TaskWorkersMapping[domain.TaskName(taskName)][domain.WorkerName(workerName)] = true
		a.repo.WorkerTasksMapping[domain.WorkerName(workerName)][domain.TaskName(taskName)] = true
	}

	wq := infra.NewWorkerQueue(workerQueueLength)
	go func(a AlphaBeeUsecase) {
	LOOP:
		for task := range a.repo.WorkerTasksMapping[domain.WorkerName(workerName)] {
			for a.repo.TaskQueues[task].Len() > 0 {
				wq <- a.repo.TaskQueues[task].Pop()

				if len(wq) == cap(wq) {
					break LOOP
				}
			}
		}
	}(a)
	a.repo.WorkerQueues[domain.WorkerName(workerName)] = wq
	return nil
}

func (a AlphaBeeUsecase) RemoveWorker(workerName string) error {
	if _, ok := a.repo.WorkerQueues[domain.WorkerName(workerName)]; !ok {
		return fmt.Errorf("worker %s not found", workerName)
	}

	// TODO: This method is very inefficient, try to find another way to
	// store worker - task mappings
	for task, _ := range a.repo.TaskWorkersMapping {
		delete(a.repo.TaskWorkersMapping[task], domain.WorkerName(workerName))
	}

	delete(a.repo.WorkerQueues, domain.WorkerName(workerName))
	return nil
}
