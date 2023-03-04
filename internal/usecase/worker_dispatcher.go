package usecase

import "context"

type WorkerDispatcher struct{}

func (d *WorkerDispatcher) Dispatch(ctx context.Context) {}
