package ktnuitygo

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/emirpasic/gods/sets/hashset"
)

func GetDefault[T any]() T {
	var zero T
	return zero
}

func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func FloatSortFunc(a float64, b float64) int {
	if math.IsNaN(a) && math.IsNaN(b) {
		return 0
	}

	if math.IsNaN(a) {
		return 1
	}

	if math.IsNaN(b) {
		return -1
	}

	if a < b {
		return -1
	}

	if a > b {
		return 1
	}

	return 0
}

func MergeSlices[T any](slices...[]T) []T {
	result := make([]T, 0, Max(len(slices), 1))
	for _, slice := range slices {
		for _, value := range slice {
			result = append(result, value)
		}
	}
	return result
}

func MergeUniqueSlices[T any](slices...[]T) []T {
	seen := hashset.New()
	result := make([]T, 0, Max(len(slices), 1))
	for _, slice := range slices {
		for _, value := range slice {
			if !seen.Contains(value) {
				seen.Add(value)
				result = append(result, value)
			}
		}
	}
	return result
}

func FirstOrDefault[T any](slice []T, or T) T {
	if len(slice) == 0 {
		return or
	}

	return slice[0]
}

func LastOrDefault[T any](slice []T, or T) T {
	if len(slice) == 0 {
		return or
	}

	return slice[len(slice) - 1]
}

type ErrorConsumerFn func(error)

func InitDefault[T any]() T {
	var zero T
	verify(&zero, false)
	return zero
}

func ForceInitDefault[T any]() T {
	var zero T
	verify(&zero, true)
	return zero
}

func verify[T any](inst *T, force bool) {
	v := reflect.ValueOf(inst).Elem()

	if v.Kind() != reflect.Struct {
		return
	}

	for i := range v.NumField() {
		field := v.Field(i)
		if !force && !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Map:
			if field.IsNil() {
				setField(field, reflect.MakeMap(field.Type()), force)
			}
		case reflect.Slice:
			if field.IsNil() {
				setField(field, reflect.MakeSlice(field.Type(), 0, 8), force)
			}
		}
	}
}

func setField(v, x reflect.Value, force bool) {
	if force {
		forceSet(v, x)
	} else {
		v.Set(x)
	}
}

// Shouldn't be used ever, our use case is for DataTank specifically, and it might use package private values that still need to be set.
func forceSet(v, x reflect.Value) {
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).
		Elem().
		Set(x)
}
