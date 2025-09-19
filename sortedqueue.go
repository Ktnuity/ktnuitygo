package ktnuitygo

import "slices"

type SortedQueue[T any] struct {
	queue []T
	sortFn func(a T, b T) int
}

func CreateSortedQueue[T any](sortFn func(a T, b T) int) *SortedQueue[T] {
	return &SortedQueue[T]{
		queue: make([]T, 0, 8),
		sortFn: sortFn,
	}
}

func (s *SortedQueue[T]) Push(data T) *SortedQueue[T] {
	s.queue = append(s.queue, data)

	return s
}

func (s *SortedQueue[T]) Raw() []T {
	result := make([]T, len(s.queue))
	copy(result, s.queue)
	return result
}

func (s *SortedQueue[T]) Sorted() []T {
	clone := s.Raw()
	slices.SortFunc(clone, s.sortFn)
	return clone
}

func (s *SortedQueue[T]) Get(index int, clamp ...bool) T {
	if len(clamp) == 0 {
		clamp = []bool{true}
	}

	sorted := s.Sorted()
	if index >= 0 {
		if clamp[0] {
			index = Min(index, len(sorted) - 1)
		}

		if index >= 0 {
			return sorted[index]
		}

		return GetDefault[T]()
	}

	index = len(sorted) + index
	if clamp[0] {
		index = Max(index, 0)
	}

	if index < len(sorted) {
		return sorted[index]
	}

	return GetDefault[T]()
}

func (s *SortedQueue[T]) Trim(offset int) []T {
	sorted := s.Sorted()
	result := make([]T, 0, len(sorted) - offset - 1)

	for i := 1; i < len(sorted) - offset; i++ {
		result = append(result, sorted[i])
	}

	return result
}
