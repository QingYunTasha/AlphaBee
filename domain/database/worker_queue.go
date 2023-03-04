package domain

type WorkerQueue struct {
	WorkerID uint
	Jobs     []Job
}
