package domain

type Algorithm string

var (
	RoundRobin       Algorithm = "ROUND_ROBIN"
	Priority         Algorithm = "PRIORITY"
	ShortestJobFirst Algorithm = "SJF"
)

type WorkerDispatcher struct {
	Algorithm Algorithm
}
