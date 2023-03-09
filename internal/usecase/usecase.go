package usecase

import (
	infradomain "AlphaBee/domain/infra"
	"AlphaBee/internal/infra"
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

	return <-wq.Jobs, nil
}

func (a AlphaBeeUsecase) AddTask(taskName string, algorithm infradomain.Algorithm, n int) error {
	if _, ok := a.repo.TaskQueues[taskName]; ok {
		return fmt.Errorf("task %s already exists", taskName)
	}

	tq := infra.NewTaskQueue(algorithm, viper.GetInt("task_queue_size"))
	a.repo.TaskQueues[taskName] = tq
	a.repo.Brokers[taskName] = infra.NewBroker(tq)

	return nil
}
func (a AlphaBeeUsecase) RemoveTask(taskName string) error {
	if _, ok := a.repo.TaskQueues[taskName]; !ok {
		return fmt.Errorf("task %s not found", taskName)
	}
	delete(a.repo.TaskQueues, taskName)
	return nil
}
func (a AlphaBeeUsecase) AddWorker(workerName string, n int) error {
	if _, ok := a.repo.WorkerQueues[workerName]; ok {
		return fmt.Errorf("worker %s already exists", workerName)
	}

	a.repo.WorkerQueues[workerName] = infradomain.WorkerQueue{
		WorkerName: workerName,
		Jobs:       make(chan infradomain.Job, n),
	}
	return nil
}
func (a AlphaBeeUsecase) RemoveWorker(workerName string) error {
	if _, ok := a.repo.WorkerQueues[workerName]; !ok {
		return fmt.Errorf("worker %s not found", workerName)
	}

	delete(a.repo.WorkerQueues, workerName)
	return nil
}
