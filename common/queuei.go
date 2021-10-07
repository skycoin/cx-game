package common

type QueueI struct {
	queue []int
}

func NewQueue() QueueI {
	queue := QueueI{
		queue: make([]int, 0),
	}
	return queue
}

func (q *QueueI) Push(n int) {
	q.queue = append(q.queue, n)
}
func (q *QueueI) Pop() int {
	if len(q.queue) == 0 {
		return -1
	}
	returnValue := q.queue[0]
	q.queue = q.queue[1:]
	return returnValue
}
