package ktnuitygo

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	content := `# Comment line
KEY1=value1
KEY2=value2
KEY3=123
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, err := LoadEnv(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load env file: %v", err)
	}

	if env.config["KEY1"] != "value1" {
		t.Errorf("Expected KEY1 to be 'value1', got '%s'", env.config["KEY1"])
	}

	if env.config["KEY2"] != "value2" {
		t.Errorf("Expected KEY2 to be 'value2', got '%s'", env.config["KEY2"])
	}

	if env.config["KEY3"] != "123" {
		t.Errorf("Expected KEY3 to be '123', got '%s'", env.config["KEY3"])
	}
}

func TestLoadEnvWithComments(t *testing.T) {
	content := `# This is a comment
KEY1=value1
// This is also a comment
KEY2=value2

/* Multi-line
comment */

KEY3=value3
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, err := LoadEnv(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load env file: %v", err)
	}

	if len(env.config) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(env.config))
	}
}

func TestLoadEnvUnclosedMultilineComment(t *testing.T) {
	content := `KEY1=value1
/*
This comment is never closed
KEY2=value2
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	_, err := LoadEnv(tmpFile)
	if err == nil {
		t.Error("Expected error for unclosed multi-line comment")
	}
}

func TestLoadEnvWithSpaces(t *testing.T) {
	content := `KEY1 = value with spaces
  KEY2  =  trimmed key
KEY3=value`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, err := LoadEnv(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load env file: %v", err)
	}

	if env.config["KEY1"] != " value with spaces" {
		t.Errorf("Expected KEY1 to preserve value spacing, got '%s'", env.config["KEY1"])
	}

	if env.config["KEY2"] != "  trimmed key" {
		t.Errorf("Expected KEY2 value to be '  trimmed key', got '%s'", env.config["KEY2"])
	}
}

func TestEnvGetString(t *testing.T) {
	content := `NAME=Alice`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	value, err := env.GetString("NAME")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if value != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", value)
	}

	_, err = env.GetString("MISSING")
	if err == nil {
		t.Error("Expected error for missing key")
	}
}

func TestEnvGetStringOrDefault(t *testing.T) {
	content := `NAME=Alice`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	value := env.GetStringOrDefault("NAME", "Bob")
	if value != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", value)
	}

	value = env.GetStringOrDefault("MISSING", "Bob")
	if value != "Bob" {
		t.Errorf("Expected default 'Bob', got '%s'", value)
	}
}

func TestEnvGetInt(t *testing.T) {
	content := `
COUNT8=127
COUNT16=32767
COUNT32=2147483647
COUNT64=9223372036854775807
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	val8, err := env.GetInt8("COUNT8")
	if err != nil || val8 != 127 {
		t.Errorf("Expected GetInt8 to return 127, got %d (err: %v)", val8, err)
	}

	val16, err := env.GetInt16("COUNT16")
	if err != nil || val16 != 32767 {
		t.Errorf("Expected GetInt16 to return 32767, got %d (err: %v)", val16, err)
	}

	val32, err := env.GetInt32("COUNT32")
	if err != nil || val32 != 2147483647 {
		t.Errorf("Expected GetInt32 to return 2147483647, got %d (err: %v)", val32, err)
	}

	val64, err := env.GetInt64("COUNT64")
	if err != nil || val64 != 9223372036854775807 {
		t.Errorf("Expected GetInt64 to return 9223372036854775807, got %d (err: %v)", val64, err)
	}
}

func TestEnvGetUint(t *testing.T) {
	content := `
UCOUNT8=255
UCOUNT16=65535
UCOUNT32=4294967295
UCOUNT64=18446744073709551615
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	val8, err := env.GetUint8("UCOUNT8")
	if err != nil || val8 != 255 {
		t.Errorf("Expected GetUint8 to return 255, got %d (err: %v)", val8, err)
	}

	val16, err := env.GetUint16("UCOUNT16")
	if err != nil || val16 != 65535 {
		t.Errorf("Expected GetUint16 to return 65535, got %d (err: %v)", val16, err)
	}

	val32, err := env.GetUint32("UCOUNT32")
	if err != nil || val32 != 4294967295 {
		t.Errorf("Expected GetUint32 to return 4294967295, got %d (err: %v)", val32, err)
	}

	val64, err := env.GetUint64("UCOUNT64")
	if err != nil || val64 != 18446744073709551615 {
		t.Errorf("Expected GetUint64 to return 18446744073709551615, got %d (err: %v)", val64, err)
	}
}

func TestEnvGetFloat(t *testing.T) {
	content := `
PI32=3.14159
PI64=3.141592653589793
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	val32, err := env.GetFloat32("PI32")
	if err != nil {
		t.Errorf("Expected no error for GetFloat32, got %v", err)
	}
	if val32 < 3.14158 || val32 > 3.14160 {
		t.Errorf("Expected GetFloat32 to return ~3.14159, got %f", val32)
	}

	val64, err := env.GetFloat64("PI64")
	if err != nil {
		t.Errorf("Expected no error for GetFloat64, got %v", err)
	}
	if val64 < 3.141592653589792 || val64 > 3.141592653589794 {
		t.Errorf("Expected GetFloat64 to return ~3.141592653589793, got %f", val64)
	}
}

func TestEnvGetBool(t *testing.T) {
	content := `
ENABLED=true
DISABLED=false
WRONG=yes
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	enabled, err := env.GetBool("ENABLED")
	if err != nil || !enabled {
		t.Errorf("Expected GetBool('ENABLED') to return true, got %v (err: %v)", enabled, err)
	}

	disabled, err := env.GetBool("DISABLED")
	if err != nil || disabled {
		t.Errorf("Expected GetBool('DISABLED') to return false, got %v (err: %v)", disabled, err)
	}

	_, err = env.GetBool("WRONG")
	if err == nil {
		t.Error("Expected error for invalid boolean value 'yes'")
	}
}

func TestEnvGetOrDefault(t *testing.T) {
	content := `COUNT=42`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	value := env.GetInt32OrDefault("COUNT", 0)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	value = env.GetInt32OrDefault("MISSING", 99)
	if value != 99 {
		t.Errorf("Expected default 99, got %d", value)
	}
}

func TestEnvHook(t *testing.T) {
	content := `ORIGINAL=value1`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	env = env.Hook(func(set EnvHookSetFn) bool {
		set("HOOKED", "hooked_value")
		set("ANOTHER", "another_value")
		return true
	})

	if env == nil {
		t.Fatal("Expected Hook to return non-nil env")
	}

	if env.config["HOOKED"] != "hooked_value" {
		t.Errorf("Expected HOOKED to be 'hooked_value', got '%s'", env.config["HOOKED"])
	}

	if env.config["ANOTHER"] != "another_value" {
		t.Errorf("Expected ANOTHER to be 'another_value', got '%s'", env.config["ANOTHER"])
	}

	// Original value should still exist
	if env.config["ORIGINAL"] != "value1" {
		t.Errorf("Expected ORIGINAL to still be 'value1', got '%s'", env.config["ORIGINAL"])
	}
}

func TestEnvHookReturnsFalse(t *testing.T) {
	content := `KEY=value`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	result := env.Hook(func(set EnvHookSetFn) bool {
		return false
	})

	if result != nil {
		t.Error("Expected Hook to return nil when hook returns false")
	}
}

func TestEnvConfig(t *testing.T) {
	content := `KEY1=value1
KEY2=value2`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	config := env.Config()
	if len(config) != 2 {
		t.Errorf("Expected config to have 2 entries, got %d", len(config))
	}

	if config["KEY1"] != "value1" {
		t.Errorf("Expected KEY1 to be 'value1', got '%s'", config["KEY1"])
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	content := `COUNT=42`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	value := GetEnvOrDefault[int32](env, "COUNT", 0)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	value = GetEnvOrDefault[int32](env, "MISSING", 99)
	if value != 99 {
		t.Errorf("Expected default 99, got %d", value)
	}
}

func TestEnvLogError(t *testing.T) {
	content := `KEY=value`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, _ := LoadEnv(tmpFile)

	errorCalled := false
	var errorFn ErrorConsumerFn = func(err error) {
		errorCalled = true
	}
	env.LogError = &errorFn

	// Trigger an error by requesting a missing key with OrDefault
	_ = env.GetStringOrDefault("MISSING", "default")

	if !errorCalled {
		t.Error("Expected LogError to be called for missing key")
	}
}

// Helper function to create temporary env files for testing
func createTempEnvFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-env-*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}
