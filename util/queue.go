package util

import (
	"container/list"
)

type Queue struct {
	list *list.List
}

func (queue *Queue) Dequeue() any {
	if queue.list.Len() > 0 {
		var item = queue.list.Front()
		queue.list.Remove(item)
		return item.Value
	}
	return nil
}

func (queue *Queue) Enqueue(value any) {
	queue.list.PushBack(value)
}

func (queue *Queue) Peek() any {
	var item = queue.list.Front()
	return item.Value
}

func (queue *Queue) Size() int {
	return queue.list.Len()
}

func NewQueue() *Queue {
	var queue = Queue{list: list.New()}
	return &queue
}
