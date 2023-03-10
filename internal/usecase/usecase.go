package usecase

import (
	infradomain "AlphaBee/domain/infra"
	"AlphaBee/internal/infra"
	taskqueue "AlphaBee/internal/infra/task_queue"
	"fmt"

	"github.com/spf13/viper"
)

type AlphaBeeUsecase struct {
	repo infradomain.Repository
}

func NewAlphaBeeUsecase(repo infradomain.Repository) *AlphaBeeUsecase {
	return &AlphaBeeUsecase{
		repo: repo,
	}
}

func (a AlphaBeeUsecase) PushJob(job infradomain.Job) error {
	a.repo.JobQueue <- job
	return nil
}

func (a AlphaBeeUsecase) PopJob(workerName string) (infradomain.Job, error) {
	wq, ok := a.repo.WorkerQueues[workerName]
	if !ok {
		return infradomain.Job{}, fmt.Errorf("worker %s not found", workerName)
	}

	return <-wq, nil
}

func (a AlphaBeeUsecase) AddTask(taskName string, algorithm infradomain.Algorithm, n int) error {
	if _, ok := a.repo.TaskQueues[taskName]; ok {
		return fmt.Errorf("task %s already exists", taskName)
	}

	tq := taskqueue.NewTaskQueue(algorithm, viper.GetInt("task_queue_size"))
	a.repo.TaskQueues[taskName] = tq
	a.repo.Brokers[taskName] = infra.NewBroker(tq)

	return nil
}

func (a AlphaBeeUsecase) RemoveTask(taskName string) error {
	if _, ok := a.repo.TaskQueues[taskName]; !ok {
		return fmt.Errorf("task %s not found", taskName)
	}
	delete(a.repo.TaskQueues, taskName)
	delete(a.repo.Brokers, taskName)
	return nil
}

func (a AlphaBeeUsecase) AddWorker(workerName string, taskNames []string, n int) error {
	if _, ok := a.repo.WorkerQueues[workerName]; ok {
		return fmt.Errorf("worker %s already exists", workerName)
	}

	// add m2m mapping
	if _, ok := a.repo.WorkerTasksMapping[workerName]; !ok {
		a.repo.WorkerTasksMapping[workerName] = make(map[string]bool)
	}

	for _, taskName := range taskNames {
		if _, ok := a.repo.TaskWorkersMapping[taskName]; !ok {
			a.repo.TaskWorkersMapping[taskName] = make(map[string]bool)
		}
		a.repo.TaskWorkersMapping[taskName][workerName] = true
		a.repo.WorkerTasksMapping[workerName][taskName] = true
	}

	// TODO: Need to initialize to fill the worker queue
	wq := infra.NewWorkerQueue(n)
	go func() {
	LOOP:
		for task := range a.repo.WorkerTasksMapping[workerName] {
			for a.repo.TaskQueues[task].Len() > 0 {
				wq <- a.repo.TaskQueues[task].Pop()

				if len(wq) == cap(wq) {
					break LOOP
				}
			}
		}
	}()
	a.repo.WorkerQueues[workerName] = wq
	return nil
}

func (a AlphaBeeUsecase) RemoveWorker(workerName string) error {
	if _, ok := a.repo.WorkerQueues[workerName]; !ok {
		return fmt.Errorf("worker %s not found", workerName)
	}

	// TODO: This method is very inefficient, try to find another way to
	// store worker - task mappings
	for _, workers := range a.repo.TaskWorkersMapping {
		delete(workers, workerName)
	}

	delete(a.repo.WorkerQueues, workerName)
	return nil
}
