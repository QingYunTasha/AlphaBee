package domain

type Job struct {
	ID       uint
	TaskName string
	Priority uint8
}

type JobQueue struct {
	Jobs chan Job
}
