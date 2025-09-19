package ktnuitygo

type LimitedQueue[T any] struct {
	items []*T
	maxSize int
}

func LimitedQueueCreate[T any](maxSize...int) *LimitedQueue[T] {
	return &LimitedQueue[T]{
		items: make([]*T, 0, FirstOrDefault(maxSize, 3)),
		maxSize: FirstOrDefault(maxSize, 3),
	}
}

func (l *LimitedQueue[T]) Peek() *T {
	if len(l.items) == 0 {
		return nil
	}

	return l.items[0]
}

func (l *LimitedQueue[T]) Push(item T) {
	if len(l.items) >= l.maxSize {
		l.items = append(l.items[1:], &item)
	} else {
		l.items = append(l.items, &item)
	}
}

func (l *LimitedQueue[T]) Pop() *T {
	if len(l.items) == 0 {
		return nil
	}

	result := l.items[0]
	l.items = l.items[1:]

	return result
}

func (l *LimitedQueue[T]) Size() int {
	return len(l.items)
}

