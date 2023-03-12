package domain

import "sync"

type Algorithm string

var (
	RoundRobin         Algorithm = "ROUND_ROBIN"
	PrioritySmallFirst Algorithm = "PRIORITY_SMALL_FIRST"
	PriorityLargeFirst Algorithm = "PRIORITY_LARGE_FIRST"
	ShortestJobFirst   Algorithm = "SHORTEST_JOB_FIRST"
)

func IsValidAlgorithm(input string) bool {
	switch Algorithm(input) {
	case RoundRobin, PriorityLargeFirst, PrioritySmallFirst, ShortestJobFirst:
		return true
	default:
		return false
	}
}

type TaskName string

type TaskQueue interface {
	Push(job Job)
	Pop() (job Job)
	Len() (n int)
}

type AsyncTaskQueue interface {
	sync.Locker
	Push(job Job)
	Pop() (job Job)
	Len() int
}
