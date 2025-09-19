package ktnuitygo

import "slices"

type LimitedArrayEntry[T comparable, V any] struct {
	value V
	data T
}

type LimitedArray[T comparable, V any] struct {
	limiter func(value V) bool
	accumulator func() V
	data []*LimitedArrayEntry[T, V]
}

func LimitedArrayCreate[T comparable, V any](limiter func(value V) bool, accumulator func() V) *LimitedArray[T, V] {
	return &LimitedArray[T, V]{
		limiter: limiter,
		accumulator: accumulator,
		data: make([]*LimitedArrayEntry[T, V], 0, 8),
	}
}

func (l *LimitedArray[T, V]) Data() []T {
	result := make([]T, len(l.data))

	for i,v := range l.data {
		result[i] = v.data
	}

	return result
}

func (l *LimitedArray[T, V]) Contains(item T) bool {
	return slices.Contains(l.Data(), item)
}

func (l *LimitedArray[T, V]) Push(item T) int {
	entry := &LimitedArrayEntry[T, V]{
		value: l.accumulator(),
		data: item,
	}

	l.data = append(l.data, entry)
	l.applyFilter()

	return len(l.data)
}

func (l *LimitedArray[T, V]) IsEmpty() bool {
	return l.Size() == 0
}

func (l *LimitedArray[T, V]) Size() int {
	return len(l.data)
}

func (l *LimitedArray[T, V]) Pop() *T {
	if len(l.data) == 0 {
		return nil
	}

	entry := l.data[len(l.data)-1]
	l.data = l.data[:len(l.data)-1]

	return &entry.data
}

func (l *LimitedArray[T, V]) applyFilter() {
	result := make([]*LimitedArrayEntry[T, V], 0, len(l.data))

	for _,v := range l.data {
		if l.limiter(v.value) {
			result = append(result, v)
		}
	}

	l.data = result
}
