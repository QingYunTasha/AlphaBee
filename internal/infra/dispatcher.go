package infra

import (
	infradomain "AlphaBee/domain/infra"
)

type Dispatcher struct {
	repo infradomain.Repository
}

func NewDispatcher(repo infradomain.Repository) infradomain.Dispatcher {
	return &Dispatcher{
		repo: repo,
	}
}

func (d Dispatcher) Run() {
	go func() {
		var job infradomain.Job
		for {
			job = <-d.repo.JobQueue

			for key, _ := range d.repo.TaskQueues {
				if job.TaskName == key {
					d.repo.TaskQueues[key].Push(job)
					break
				}
			}
		}
	}()
}
