package domain

type Algorithm string

var (
	RoundRobin       Algorithm = "ROUND_ROBIN"
	Priority         Algorithm = "PRIORITY"
	ShortestJobFirst Algorithm = "SJF"
)

type Task struct {
	ID         uint
	Name       string
	Preemption bool
	Algorithm  Algorithm
}
