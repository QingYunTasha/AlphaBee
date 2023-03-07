package domain

type Algorithm string

var (
	RoundRobin       Algorithm = "ROUND_ROBIN"
	PriorityFirst    Algorithm = "PRIORITY_FIRST"
	ShortestJobFirst Algorithm = "SHORTEST_JOB_FIRST"
)

type TaskQueue struct {
	ID        uint
	TaskName  string
	Algorithm Algorithm
	Jobs      []Job
}
