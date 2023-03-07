package utils

type Node struct {
	Priority uint
}

type PriorityQueue []Node

func NewPriorityQueue() PriorityQueue { return append(PriorityQueue{}, Node{}) }

func (q PriorityQueue) less(i, j int) bool { return q[i].Priority > q[j].Priority }

func (q PriorityQueue) swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q *PriorityQueue) Push(n Node) {
	*q = append(*q, n)
	q.up(len(*q) - 1)
}

func (q *PriorityQueue) Pop() Node {
	n := len(*q) - 1
	res := (*q)[1]
	(*q)[1] = (*q)[n]
	*q = (*q)[:n]
	q.down(1)
	return res
}

func (q PriorityQueue) up(i int) {
	for j := i >> 1; i > 1 && q.less(j, i); i, j = j, j>>1 {
		q.swap(i, j)
	}
}

func (q PriorityQueue) down(i int) {
	for j := i << 1; j < len(q); i, j = j, j<<1 {
		if j+1 < len(q) && q.less(j, j+1) {
			j++
		}
		if q.less(j, i) {
			break
		}
		q.swap(i, j)
	}
}
