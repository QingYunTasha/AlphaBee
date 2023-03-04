package domain

type TaskQueue struct {
	TaskID uint
	Jobs   []Job
}
