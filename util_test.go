package ktnuitygo

import (
	"math"
	"testing"
)

func TestGetDefault(t *testing.T) {
	intVal := GetDefault[int]()
	if intVal != 0 {
		t.Errorf("Expected GetDefault[int]() to be 0, got %d", intVal)
	}

	strVal := GetDefault[string]()
	if strVal != "" {
		t.Errorf("Expected GetDefault[string]() to be empty string, got '%s'", strVal)
	}

	boolVal := GetDefault[bool]()
	if boolVal != false {
		t.Errorf("Expected GetDefault[bool]() to be false, got %v", boolVal)
	}
}

func TestMin(t *testing.T) {
	if Min(5, 10) != 5 {
		t.Errorf("Expected Min(5, 10) to be 5, got %d", Min(5, 10))
	}

	if Min(10, 5) != 5 {
		t.Errorf("Expected Min(10, 5) to be 5, got %d", Min(10, 5))
	}

	if Min(3.5, 2.1) != 2.1 {
		t.Errorf("Expected Min(3.5, 2.1) to be 2.1, got %f", Min(3.5, 2.1))
	}

	if Min(-5, -10) != -10 {
		t.Errorf("Expected Min(-5, -10) to be -10, got %d", Min(-5, -10))
	}
}

func TestMax(t *testing.T) {
	if Max(5, 10) != 10 {
		t.Errorf("Expected Max(5, 10) to be 10, got %d", Max(5, 10))
	}

	if Max(10, 5) != 10 {
		t.Errorf("Expected Max(10, 5) to be 10, got %d", Max(10, 5))
	}

	if Max(3.5, 2.1) != 3.5 {
		t.Errorf("Expected Max(3.5, 2.1) to be 3.5, got %f", Max(3.5, 2.1))
	}

	if Max(-5, -10) != -5 {
		t.Errorf("Expected Max(-5, -10) to be -5, got %d", Max(-5, -10))
	}
}

func TestFloatSortFunc(t *testing.T) {
	// Test normal comparisons
	if FloatSortFunc(1.0, 2.0) != -1 {
		t.Error("Expected FloatSortFunc(1.0, 2.0) to return -1")
	}

	if FloatSortFunc(2.0, 1.0) != 1 {
		t.Error("Expected FloatSortFunc(2.0, 1.0) to return 1")
	}

	if FloatSortFunc(1.0, 1.0) != 0 {
		t.Error("Expected FloatSortFunc(1.0, 1.0) to return 0")
	}

	// Test NaN handling
	nan := math.NaN()

	if FloatSortFunc(nan, nan) != 0 {
		t.Error("Expected FloatSortFunc(NaN, NaN) to return 0")
	}

	if FloatSortFunc(nan, 1.0) != 1 {
		t.Error("Expected FloatSortFunc(NaN, 1.0) to return 1")
	}

	if FloatSortFunc(1.0, nan) != -1 {
		t.Error("Expected FloatSortFunc(1.0, NaN) to return -1")
	}
}

func TestMergeSlices(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	slice3 := []int{7, 8}

	result := MergeSlices(slice1, slice2, slice3)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}

	if len(result) != len(expected) {
		t.Errorf("Expected merged length to be %d, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestMergeSlicesEmpty(t *testing.T) {
	result := MergeSlices([]int{}, []int{}, []int{})
	if len(result) != 0 {
		t.Errorf("Expected empty merge to have length 0, got %d", len(result))
	}
}

func TestMergeUniqueSlices(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{2, 3, 4}
	slice3 := []int{3, 4, 5}

	result := MergeUniqueSlices(slice1, slice2, slice3)
	expected := []int{1, 2, 3, 4, 5}

	if len(result) != len(expected) {
		t.Errorf("Expected unique merged length to be %d, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("At index %d, expected %d, got %d", i, expected[i], result[i])
		}
	}
}

func TestMergeUniqueSlicesOrder(t *testing.T) {
	// Verify that order is preserved (first occurrence)
	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"b", "c", "d"}

	result := MergeUniqueSlices(slice1, slice2)
	expected := []string{"a", "b", "c", "d"}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("At index %d, expected '%s', got '%s'", i, expected[i], result[i])
		}
	}
}

func TestFirstOrDefault(t *testing.T) {
	slice := []int{10, 20, 30}
	result := FirstOrDefault(slice, 99)
	if result != 10 {
		t.Errorf("Expected FirstOrDefault to return 10, got %d", result)
	}

	emptySlice := []int{}
	result = FirstOrDefault(emptySlice, 99)
	if result != 99 {
		t.Errorf("Expected FirstOrDefault on empty slice to return 99, got %d", result)
	}
}

func TestLastOrDefault(t *testing.T) {
	slice := []int{10, 20, 30}
	result := LastOrDefault(slice, 99)
	if result != 30 {
		t.Errorf("Expected LastOrDefault to return 30, got %d", result)
	}

	emptySlice := []int{}
	result = LastOrDefault(emptySlice, 99)
	if result != 99 {
		t.Errorf("Expected LastOrDefault on empty slice to return 99, got %d", result)
	}
}

func TestInitDefault(t *testing.T) {
	type TestStruct struct {
		Items []int
		Meta  map[string]string
	}

	result := InitDefault[TestStruct]()

	if result.Items == nil {
		t.Error("Expected Items slice to be initialized")
	}

	if result.Meta == nil {
		t.Error("Expected Meta map to be initialized")
	}

	if len(result.Items) != 0 {
		t.Errorf("Expected Items to be empty, got length %d", len(result.Items))
	}

	if len(result.Meta) != 0 {
		t.Errorf("Expected Meta to be empty, got length %d", len(result.Meta))
	}
}

func TestForceInitDefault(t *testing.T) {
	type TestStruct struct {
		items []int // private field
		meta  map[string]string // private field
	}

	result := ForceInitDefault[TestStruct]()

	if result.items == nil {
		t.Error("Expected items slice to be force-initialized")
	}

	if result.meta == nil {
		t.Error("Expected meta map to be force-initialized")
	}
}

func TestInitDefaultWithNestedStructs(t *testing.T) {
	type Inner struct {
		Data []string
	}

	type Outer struct {
		Inner Inner
		Items []int
	}

	result := InitDefault[Outer]()

	if result.Items == nil {
		t.Error("Expected Items slice to be initialized")
	}

	// Note: nested struct fields are not automatically initialized by verify
	// This is expected behavior
}
