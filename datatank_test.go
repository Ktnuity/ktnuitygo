package ktnuitygo

import (
	"os"
	"testing"
)

type TestData struct {
	Name  string
	Count int
	Items []string
	Meta  map[string]int
}

type TestDataObject struct {
	Name		string		`json:"name"`
	Age			int			`json:"age"`
}

type TestDataComplex struct {
	Data	map[string]TestDataObject
}

func TestDataTankSetDir(t *testing.T) {
	DataTankSetDir("/tmp/test")
	if tankDir != "/tmp/test" {
		t.Errorf("Expected tankDir to be '/tmp/test', got '%s'", tankDir)
	}

	DataTankSetDir("/tmp/test/")
	if tankDir != "/tmp/test" {
		t.Errorf("Expected tankDir to be '/tmp/test' (trailing slash removed), got '%s'", tankDir)
	}

	DataTankSetDir("")
	if tankDir != "." {
		t.Errorf("Expected tankDir to be '.', got '%s'", tankDir)
	}
}

func TestDataTankNew(t *testing.T) {
	tmpDir := t.TempDir()
	DataTankSetDir(tmpDir)

	tank, err := DataTankNew[TestData]("test")
	if err != nil {
		t.Fatalf("Failed to create DataTank: %v", err)
	}

	if tank.name != "test" {
		t.Errorf("Expected tank name to be 'test', got '%s'", tank.name)
	}

	if tank.data == nil {
		t.Fatal("Expected tank data to be initialized")
	}

	// Verify that slices and maps are initialized
	if tank.data.Items == nil {
		t.Error("Expected Items slice to be initialized")
	}

	if tank.data.Meta == nil {
		t.Error("Expected Meta map to be initialized")
	}
}

func TestDataTankSaveAndReload(t *testing.T) {
	tmpDir := t.TempDir()
	DataTankSetDir(tmpDir)

	tank, err := DataTankNew[TestData]("test-save")
	if err != nil {
		t.Fatalf("Failed to create DataTank: %v", err)
	}

	// Set some data
	err = DataTankSet(tank, func(data *TestData) {
		data.Name = "Alice"
		data.Count = 42
		data.Items = []string{"foo", "bar", "baz"}
		data.Meta = map[string]int{"x": 1, "y": 2}
	})
	if err != nil {
		t.Fatalf("Failed to set DataTank data: %v", err)
	}

	// Verify data was set
	result := DataTankGet(tank, func(data *TestData) *string {
		return &data.Name
	})
	if *result != "Alice" {
		t.Errorf("Expected Name to be 'Alice', got '%s'", *result)
	}

	// Reload and verify
	err = tank.Reload()
	if err != nil {
		t.Fatalf("Failed to reload DataTank: %v", err)
	}

	if tank.data.Name != "Alice" {
		t.Errorf("After reload, expected Name to be 'Alice', got '%s'", tank.data.Name)
	}
	if tank.data.Count != 42 {
		t.Errorf("After reload, expected Count to be 42, got %d", tank.data.Count)
	}
	if len(tank.data.Items) != 3 {
		t.Errorf("After reload, expected 3 items, got %d", len(tank.data.Items))
	}
	if tank.data.Meta["x"] != 1 {
		t.Errorf("After reload, expected Meta['x'] to be 1, got %d", tank.data.Meta["x"])
	}
}

func TestDataTankSaveAndReloadComplex(t *testing.T) {
	tmpDir := t.TempDir()
	DataTankSetDir(tmpDir)

	tank, err := DataTankNew[TestDataComplex]("test-data-complex")
	if err != nil {
		t.Fatalf("Failed to create DataTank: %v", err)
	}

	// Set some data
	err = DataTankSet(tank, func(data *TestDataComplex) {
		data.Data["Alice"] = TestDataObject{
			Name: "Alice",
			Age: 42,
		}
	})
	if err != nil {
		t.Fatalf("Failed to set DataTank data: %v", err)
	}

	err = DataTankSet(tank, func(data *TestDataComplex) {
		data.Data["Bob"] = TestDataObject{
			Name: "Bob",
			Age: 32,
		}
	})
	if err != nil {
		t.Fatalf("Failed to set DataTank data: %v", err)
	}

	// Verify data was set
	result := DataTankGet(tank, func(data *TestDataComplex) *TestDataObject {
		result, exists := data.Data["Alice"]
		if !exists { return nil }
		return &result
	})
	if result == nil {
		t.Errorf("Expected Data to contain Alice object.")
	}

	if result.Name != "Alice" {
		t.Errorf("Expected Name to be 'Alice', got '%s'", result.Name)
	}

	if result.Age != 42 {
		t.Errorf("Expected Age to be 42, got %d", result.Age)
	}

	// Verify missing is not set
	result = DataTankGet(tank, func(data *TestDataComplex) *TestDataObject {
		result, exists := data.Data["MISSING"]
		if !exists { return nil }
		return &result
	})
	if result != nil {
		t.Errorf("Expected Data to be nil, got %v", *result)
	}

	// Reload and verify
	err = tank.Reload()
	if err != nil {
		t.Fatalf("Failed to reload DataTank: %v", err)
	}

	name := DataTankGet(tank, func(data *TestDataComplex) *string {
		result, exists := data.Data["Alice"]
		if !exists { return nil }
		return &result.Name
	})
	if name == nil {
		t.Errorf("After reload, expected Name to be non-nil")
	}
	if *name != "Alice" {
		t.Errorf("After reload, expected Name to be 'Alice', got '%s'", *name)
	}

	name = DataTankGet(tank, func(data *TestDataComplex) *string {
		result, exists := data.Data["MISSING"]
		if !exists { return nil }
		return &result.Name
	})
	if name != nil {
		t.Errorf("After reload, expected Name to be nil, got '%s'", *name)
	}

	// Re-create and verify
	tank, err = DataTankNew[TestDataComplex]("test-data-complex")
	if err != nil {
		t.Fatalf("Failed to create DataTank: %v", err)
	}

	name = DataTankGet(tank, func(data *TestDataComplex) *string {
		result, exists := data.Data["Alice"]
		if !exists { return nil }
		return &result.Name
	})

	if name == nil {
		t.Errorf("After recreate, expected Name to be non-nil")
	}
	if *name != "Alice" {
		t.Errorf("After recreate, expected Name to be 'Alice', got '%s'", *name)
	}
}

func TestDataTankGet(t *testing.T) {
	tmpDir := t.TempDir()
	DataTankSetDir(tmpDir)

	tank, err := DataTankNew[TestData]("test-get")
	if err != nil {
		t.Fatalf("Failed to create DataTank: %v", err)
	}

	DataTankSet(tank, func(data *TestData) {
		data.Count = 100
	})

	count := DataTankGet(tank, func(data *TestData) *int {
		return &data.Count
	})

	if *count != 100 {
		t.Errorf("Expected count to be 100, got %d", *count)
	}
}

func TestTankLoadError(t *testing.T) {
	tmpDir := t.TempDir()
	DataTankSetDir(tmpDir)

	// Create a corrupted JSON file
	badFile := tmpDir + "/bad.tank.json"
	err := os.WriteFile(badFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create bad file: %v", err)
	}

	_, err = DataTankNew[TestData]("bad")
	if err == nil {
		t.Error("Expected error when loading corrupted JSON")
	}
}

func TestTankPath(t *testing.T) {
	DataTankSetDir("/test/dir")
	path := tankPath("mydata")
	expected := "/test/dir/mydata.tank.json"
	if path != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, path)
	}
}

