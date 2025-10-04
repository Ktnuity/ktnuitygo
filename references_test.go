package ktnuitygo

import "testing"

func TestAsRef(t *testing.T) {
	value := 42
	ref := AsRef(value)

	if ref == nil {
		t.Fatal("Expected AsRef to return non-nil pointer")
	}

	if *ref != 42 {
		t.Errorf("Expected referenced value to be 42, got %d", *ref)
	}

	// Verify it's a copy, not the original
	value = 100
	if *ref == 100 {
		t.Error("Expected AsRef to create a copy, but it references the original")
	}
}

func TestAsRefString(t *testing.T) {
	str := "hello"
	ref := AsRef(str)

	if ref == nil {
		t.Fatal("Expected AsRef to return non-nil pointer")
	}

	if *ref != "hello" {
		t.Errorf("Expected referenced value to be 'hello', got '%s'", *ref)
	}
}

func TestAsRefStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Alice", Age: 30}
	ref := AsRef(person)

	if ref == nil {
		t.Fatal("Expected AsRef to return non-nil pointer")
	}

	if ref.Name != "Alice" || ref.Age != 30 {
		t.Errorf("Expected referenced person to be {Alice, 30}, got {%s, %d}", ref.Name, ref.Age)
	}
}

func TestAsRefMany(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	refs := AsRefMany(values)

	if len(refs) != len(values) {
		t.Errorf("Expected %d references, got %d", len(values), len(refs))
	}

	for i, ref := range refs {
		if ref == nil {
			t.Errorf("Expected reference at index %d to be non-nil", i)
			continue
		}

		if *ref != values[i] {
			t.Errorf("At index %d, expected referenced value to be %d, got %d", i, values[i], *ref)
		}
	}
}

func TestAsRefManyEmpty(t *testing.T) {
	values := []int{}
	refs := AsRefMany(values)

	if len(refs) != 0 {
		t.Errorf("Expected empty slice result, got length %d", len(refs))
	}
}

func TestAsRefManyStrings(t *testing.T) {
	values := []string{"foo", "bar", "baz"}
	refs := AsRefMany(values)

	if len(refs) != 3 {
		t.Errorf("Expected 3 references, got %d", len(refs))
	}

	for i, ref := range refs {
		if *ref != values[i] {
			t.Errorf("At index %d, expected '%s', got '%s'", i, values[i], *ref)
		}
	}
}

func TestAsRefManyCopies(t *testing.T) {
	values := []int{10, 20, 30}
	refs := AsRefMany(values)

	// Modify original slice
	values[0] = 999

	// Verify references still point to original values
	if *refs[0] == 999 {
		t.Error("Expected AsRefMany to create copies, but it references the original values")
	}

	if *refs[0] != 10 {
		t.Errorf("Expected first reference to be 10, got %d", *refs[0])
	}
}
