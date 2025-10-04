package ktnuitygo

import "testing"

func TestLimitedQueueCreate(t *testing.T) {
	queue := LimitedQueueCreate[int]()
	if queue.maxSize != 3 {
		t.Errorf("Expected default maxSize to be 3, got %d", queue.maxSize)
	}

	queue = LimitedQueueCreate[int](5)
	if queue.maxSize != 5 {
		t.Errorf("Expected maxSize to be 5, got %d", queue.maxSize)
	}
}

func TestLimitedQueuePushAndPeek(t *testing.T) {
	queue := LimitedQueueCreate[int](3)

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	peeked := queue.Peek()
	if peeked == nil || *peeked != 1 {
		t.Errorf("Expected Peek to return 1, got %v", peeked)
	}

	if queue.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", queue.Size())
	}
}

func TestLimitedQueuePushExceedLimit(t *testing.T) {
	queue := LimitedQueueCreate[int](3)

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)
	queue.Push(4) // This should remove 1

	peeked := queue.Peek()
	if peeked == nil || *peeked != 2 {
		t.Errorf("Expected Peek to return 2 after exceeding limit, got %v", peeked)
	}

	if queue.Size() != 3 {
		t.Errorf("Expected size to remain 3, got %d", queue.Size())
	}

	queue.Push(5) // This should remove 2
	peeked = queue.Peek()
	if peeked == nil || *peeked != 3 {
		t.Errorf("Expected Peek to return 3, got %v", peeked)
	}
}

func TestLimitedQueuePop(t *testing.T) {
	queue := LimitedQueueCreate[string](3)

	queue.Push("a")
	queue.Push("b")
	queue.Push("c")

	popped := queue.Pop()
	if popped == nil || *popped != "a" {
		t.Errorf("Expected Pop to return 'a', got %v", popped)
	}

	if queue.Size() != 2 {
		t.Errorf("Expected size to be 2 after pop, got %d", queue.Size())
	}

	popped = queue.Pop()
	if popped == nil || *popped != "b" {
		t.Errorf("Expected Pop to return 'b', got %v", popped)
	}

	popped = queue.Pop()
	if popped == nil || *popped != "c" {
		t.Errorf("Expected Pop to return 'c', got %v", popped)
	}

	// Pop from empty queue
	popped = queue.Pop()
	if popped != nil {
		t.Errorf("Expected Pop on empty queue to return nil, got %v", popped)
	}
}

func TestLimitedQueuePeekEmpty(t *testing.T) {
	queue := LimitedQueueCreate[int]()

	peeked := queue.Peek()
	if peeked != nil {
		t.Errorf("Expected Peek on empty queue to return nil, got %v", peeked)
	}
}

func TestLimitedQueueSize(t *testing.T) {
	queue := LimitedQueueCreate[int](5)

	if queue.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", queue.Size())
	}

	queue.Push(1)
	queue.Push(2)

	if queue.Size() != 2 {
		t.Errorf("Expected size to be 2, got %d", queue.Size())
	}

	queue.Pop()

	if queue.Size() != 1 {
		t.Errorf("Expected size to be 1 after pop, got %d", queue.Size())
	}
}
