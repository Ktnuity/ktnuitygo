package ktnuitygo

import "testing"

func TestLimitedArrayCreate(t *testing.T) {
	limiter := func(v int) bool { return v < 10 }
	accumulator := func() int { return 0 }

	arr := LimitedArrayCreate[string, int](limiter, accumulator)

	if arr == nil {
		t.Fatal("Expected LimitedArrayCreate to return non-nil")
	}

	if arr.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", arr.Size())
	}
}

func TestLimitedArrayPushAndData(t *testing.T) {
	counter := 0
	limiter := func(v int) bool { return v < 5 }
	accumulator := func() int {
		counter++
		return counter
	}

	arr := LimitedArrayCreate[string, int](limiter, accumulator)

	arr.Push("first")
	arr.Push("second")
	arr.Push("third")

	data := arr.Data()
	if len(data) != 3 {
		t.Errorf("Expected data length to be 3, got %d", len(data))
	}

	if data[0] != "first" || data[1] != "second" || data[2] != "third" {
		t.Errorf("Expected data to be ['first', 'second', 'third'], got %v", data)
	}
}

func TestLimitedArrayPushWithFilter(t *testing.T) {
	counter := 0
	limiter := func(v int) bool { return v <= 3 } // Only keep items with counter <= 3
	accumulator := func() int {
		counter++
		return counter
	}

	arr := LimitedArrayCreate[string, int](limiter, accumulator)

	arr.Push("first")   // counter = 1, kept
	arr.Push("second")  // counter = 2, kept
	arr.Push("third")   // counter = 3, kept
	arr.Push("fourth")  // counter = 4, removed by filter

	if arr.Size() != 3 {
		t.Errorf("Expected size to be 3 after filter, got %d", arr.Size())
	}

	data := arr.Data()
	if len(data) != 3 {
		t.Errorf("Expected data length to be 3, got %d", len(data))
	}

	// Verify the fourth item was filtered out
	if arr.Contains("fourth") {
		t.Error("Expected 'fourth' to be filtered out")
	}
}

func TestLimitedArrayContains(t *testing.T) {
	limiter := func(v int) bool { return true }
	accumulator := func() int { return 0 }

	arr := LimitedArrayCreate[int, int](limiter, accumulator)

	arr.Push(10)
	arr.Push(20)
	arr.Push(30)

	if !arr.Contains(20) {
		t.Error("Expected array to contain 20")
	}

	if arr.Contains(40) {
		t.Error("Expected array to not contain 40")
	}
}

func TestLimitedArrayIsEmpty(t *testing.T) {
	limiter := func(v int) bool { return true }
	accumulator := func() int { return 0 }

	arr := LimitedArrayCreate[string, int](limiter, accumulator)

	if !arr.IsEmpty() {
		t.Error("Expected new array to be empty")
	}

	arr.Push("item")

	if arr.IsEmpty() {
		t.Error("Expected array to not be empty after push")
	}
}

func TestLimitedArrayPop(t *testing.T) {
	limiter := func(v int) bool { return true }
	accumulator := func() int { return 0 }

	arr := LimitedArrayCreate[string, int](limiter, accumulator)

	arr.Push("first")
	arr.Push("second")
	arr.Push("third")

	popped := arr.Pop()
	if popped == nil || *popped != "third" {
		t.Errorf("Expected Pop to return 'third', got %v", popped)
	}

	if arr.Size() != 2 {
		t.Errorf("Expected size to be 2 after pop, got %d", arr.Size())
	}

	popped = arr.Pop()
	if popped == nil || *popped != "second" {
		t.Errorf("Expected Pop to return 'second', got %v", popped)
	}

	popped = arr.Pop()
	if popped == nil || *popped != "first" {
		t.Errorf("Expected Pop to return 'first', got %v", popped)
	}

	// Pop from empty array
	popped = arr.Pop()
	if popped != nil {
		t.Errorf("Expected Pop on empty array to return nil, got %v", popped)
	}
}

func TestLimitedArraySize(t *testing.T) {
	limiter := func(v int) bool { return true }
	accumulator := func() int { return 0 }

	arr := LimitedArrayCreate[int, int](limiter, accumulator)

	if arr.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", arr.Size())
	}

	result := arr.Push(1)
	if result != 1 {
		t.Errorf("Expected Push to return size 1, got %d", result)
	}

	arr.Push(2)
	arr.Push(3)

	if arr.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", arr.Size())
	}
}
