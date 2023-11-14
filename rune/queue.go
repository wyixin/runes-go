package rune

import "errors"

type opcode struct {
	Code byte
	Data []byte
}

type queue []*opcode

func newQueue() *queue {
	return &queue{}
}

func (q *queue) Enqueue(item *opcode) {
	*q = append(*q, item)
}

func (q *queue) Dequeue() (*opcode, error) {
	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}

	item := (*q)[0]
	*q = (*q)[1:]
	return item, nil
}

func (q *queue) IsEmpty() bool {
	return len(*q) == 0
}
