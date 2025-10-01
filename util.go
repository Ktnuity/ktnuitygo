package ktnuitygo

import (
	"math"
	"reflect"

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
	v := reflect.ValueOf(&zero).Elem()

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.CanSet() {
				continue
			}

			switch field.Kind() {
			case reflect.Map:
				field.Set(reflect.MakeMap(field.Type()))
			case reflect.Slice:
				field.Set(reflect.MakeSlice(field.Type(), 0, 8))
			}
		}
	}

	return zero
}
