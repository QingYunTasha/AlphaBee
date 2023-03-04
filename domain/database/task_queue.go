package domain

type Algorithm string

var (
	RoundRobin       Algorithm = "ROUND_ROBIN"
	Priority         Algorithm = "PRIORITY"
	ShortestJobFirst Algorithm = "SJF"
)

type TaskQueue struct {
	TaskID     uint
	Preemption bool
	Algorithm  Algorithm
	Jobs       []Job
}
