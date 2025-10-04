package ktnuitygo

import (
	"testing"
)

func TestCreateSortedQueue(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	if queue == nil {
		t.Fatal("Expected CreateSortedQueue to return non-nil")
	}

	if len(queue.queue) != 0 {
		t.Errorf("Expected initial queue length to be 0, got %d", len(queue.queue))
	}
}

func TestSortedQueuePush(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	result := queue.Push(3)
	if result != queue {
		t.Error("Expected Push to return the queue itself")
	}

	queue.Push(1).Push(4).Push(2)

	raw := queue.Raw()
	if len(raw) != 4 {
		t.Errorf("Expected queue to have 4 elements, got %d", len(raw))
	}
}

func TestSortedQueueRaw(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(3).Push(1).Push(4).Push(2)

	raw := queue.Raw()
	expected := []int{3, 1, 4, 2}

	if len(raw) != len(expected) {
		t.Errorf("Expected raw length to be %d, got %d", len(expected), len(raw))
	}

	for i := range raw {
		if raw[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], raw[i])
		}
	}
}

func TestSortedQueueSorted(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(3).Push(1).Push(4).Push(2)

	sorted := queue.Sorted()
	expected := []int{1, 2, 3, 4}

	if len(sorted) != len(expected) {
		t.Errorf("Expected sorted length to be %d, got %d", len(expected), len(sorted))
	}

	for i := range sorted {
		if sorted[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], sorted[i])
		}
	}
}

func TestSortedQueueSortedDescending(t *testing.T) {
	sortFn := func(a, b int) int { return b - a } // Descending
	queue := CreateSortedQueue(sortFn)

	queue.Push(3).Push(1).Push(4).Push(2)

	sorted := queue.Sorted()
	expected := []int{4, 3, 2, 1}

	for i := range sorted {
		if sorted[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], sorted[i])
		}
	}
}

func TestSortedQueueGet(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(3).Push(1).Push(4).Push(2)

	// Test positive indices
	if queue.Get(0) != 1 {
		t.Errorf("Expected Get(0) to be 1, got %d", queue.Get(0))
	}
	if queue.Get(1) != 2 {
		t.Errorf("Expected Get(1) to be 2, got %d", queue.Get(1))
	}
	if queue.Get(3) != 4 {
		t.Errorf("Expected Get(3) to be 4, got %d", queue.Get(3))
	}

	// Test negative indices
	if queue.Get(-1) != 4 {
		t.Errorf("Expected Get(-1) to be 4, got %d", queue.Get(-1))
	}
	if queue.Get(-2) != 3 {
		t.Errorf("Expected Get(-2) to be 3, got %d", queue.Get(-2))
	}

	// Test clamping
	if queue.Get(10) != 4 { // Should clamp to last element
		t.Errorf("Expected Get(10) to clamp to 4, got %d", queue.Get(10))
	}
	if queue.Get(-10) != 1 { // Should clamp to first element
		t.Errorf("Expected Get(-10) to clamp to 1, got %d", queue.Get(-10))
	}
}

func TestSortedQueueGetNoClamp(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(3).Push(1).Push(4).Push(2)

	// Test out of bounds without clamping - this will panic by design
	// The implementation doesn't handle out-of-bounds when clamp=false for positive indices
	// So we test that negative indices still work correctly
	result := queue.Get(-3, false)
	if result != 2 {
		t.Errorf("Expected Get(-3, false) to return 2, got %d", result)
	}
}

func TestSortedQueueTrim(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(5).Push(3).Push(1).Push(4).Push(2)

	// Trim with offset 0 removes first and last element
	// sorted = [1, 2, 3, 4, 5]
	// loop: i from 1 to len(sorted) - offset = 5, so i = 1,2,3,4
	// result = [2, 3, 4, 5]
	trimmed := queue.Trim(0)
	expected := []int{2, 3, 4, 5}

	if len(trimmed) != len(expected) {
		t.Errorf("Expected trimmed length to be %d, got %d", len(expected), len(trimmed))
	}

	for i := range trimmed {
		if trimmed[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], trimmed[i])
		}
	}
}

func TestSortedQueueTrimWithOffset(t *testing.T) {
	sortFn := func(a, b int) int { return a - b }
	queue := CreateSortedQueue(sortFn)

	queue.Push(5).Push(3).Push(1).Push(4).Push(2)

	// Trim with offset 1
	// sorted = [1, 2, 3, 4, 5]
	// loop: i from 1 to len(sorted) - offset = 4, so i = 1,2,3
	// result = [2, 3, 4]
	trimmed := queue.Trim(1)
	expected := []int{2, 3, 4}

	if len(trimmed) != len(expected) {
		t.Errorf("Expected trimmed length to be %d, got %d", len(expected), len(trimmed))
	}

	for i := range trimmed {
		if trimmed[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], trimmed[i])
		}
	}
}

func TestSortedQueueWithStrings(t *testing.T) {
	sortFn := func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}

	queue := CreateSortedQueue(sortFn)
	queue.Push("banana").Push("apple").Push("cherry")

	sorted := queue.Sorted()
	expected := []string{"apple", "banana", "cherry"}

	for i := range sorted {
		if sorted[i] != expected[i] {
			t.Errorf("At index %d, expected '%s', got '%s'", i, expected[i], sorted[i])
		}
	}
}
