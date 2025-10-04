package ktnuitygo

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"
)

type EnvData struct {
	config map[string]string
	LogError *ErrorConsumerFn
}

func LoadEnv(path...string) (*EnvData, error) {
	filepath := "./.env"
	if len(path) != 0 {
		filepath = path[0]
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load env file '%s': %v", filepath, err)
	}

	defer file.Close()
	config := make(map[string]string)
	scanner := bufio.NewScanner(file)

	multiLineComment := false

	for scanner.Scan() {
		line := scanner.Text()
		blankLine := strings.TrimSpace(line)

		if blankLine == "" ||
			strings.HasPrefix(blankLine, "#") ||
			strings.HasPrefix(blankLine, "//") {
			continue
		}

		if multiLineComment {
			if blankLine == "*/" {
				multiLineComment = false
			}

			continue
		} else if blankLine == "/*" {
			multiLineComment = true
			continue
		}

		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[: idx])
			value := line[idx + 1 :]
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if multiLineComment {
		return nil, fmt.Errorf("cannot end on a multi-line comment")
	}

	return &EnvData{
		config: config,
	}, nil
}

func consume[T any](e *EnvData, value T, err error, or T) T {
	if err != nil {
		if e.LogError != nil {
			(*e.LogError)(err)
		}

		return or
	}

	return value
}

type EnvHookSetFn func(name string, value string)
type EnvHookFn func(set EnvHookSetFn) bool

func (e *EnvData) Hook(hook EnvHookFn) *EnvData {
	config := make(map[string]string)
	setFn := func(name string, value string) {
		config[name] = value
	}

	if !hook(setFn) {
		return nil
	}

	maps.Copy(e.config, config)
	return e
}

func (e *EnvData) Config() map[string]string {
	return e.config
}

func (e *EnvData) GetString(name string) (string, error) {
	value, exists := e.config[name]
	if !exists {
		return "", fmt.Errorf("env key '%s' not found", name)
	}

	return value, nil
}

func (e *EnvData) GetStringOrDefault(name string, or string) string {
	value, err := e.GetString(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetInt8(name string) (int8, error) {
	return GetEnv[int8](e, name)
}

func (e *EnvData) GetInt8OrDefault(name string, or int8) int8 {
	value, err := e.GetInt8(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetUint8(name string) (uint8, error) {
	return GetEnv[uint8](e, name)
}

func (e *EnvData) GetUint8OrDefault(name string, or uint8) uint8 {
	value, err := e.GetUint8(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetInt16(name string) (int16, error) {
	return GetEnv[int16](e, name)
}

func (e *EnvData) GetInt16OrDefault(name string, or int16) int16 {
	value, err := e.GetInt16(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetUint16(name string) (uint16, error) {
	return GetEnv[uint16](e, name)
}

func (e *EnvData) GetUint16OrDefault(name string, or uint16) uint16 {
	value, err := e.GetUint16(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetInt32(name string) (int32, error) {
	return GetEnv[int32](e, name)
}

func (e *EnvData) GetInt32OrDefault(name string, or int32) int32 {
	value, err := e.GetInt32(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetUint32(name string) (uint32, error) {
	return GetEnv[uint32](e, name)
}

func (e *EnvData) GetUint32OrDefault(name string, or uint32) uint32 {
	value, err := e.GetUint32(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetInt64(name string) (int64, error) {
	return GetEnv[int64](e, name)
}

func (e *EnvData) GetInt64OrDefault(name string, or int64) int64 {
	value, err := e.GetInt64(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetUint64(name string) (uint64, error) {
	return GetEnv[uint64](e, name)
}

func (e *EnvData) GetUint64OrDefault(name string, or uint64) uint64 {
	value, err := e.GetUint64(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetFloat32(name string) (float32, error) {
	return GetEnv[float32](e, name)
}

func (e *EnvData) GetFloat32OrDefault(name string, or float32) float32 {
	value, err := e.GetFloat32(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetFloat64(name string) (float64, error) {
	return GetEnv[float64](e, name)
}

func (e *EnvData) GetFloat64OrDefault(name string, or float64) float64 {
	value, err := e.GetFloat64(name)
	return consume(e, value, err, or)
}

func (e *EnvData) GetBool(name string) (bool, error) {
	return GetEnv[bool](e, name)
}

func (e *EnvData) GetBoolOrDefault(name string, or bool) bool {
	value, err := e.GetBool(name)
	return consume(e, value, err, or)
}

type EnvValueType interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~float32 | ~float64 | ~bool | ~string
}

func GetEnvOrDefault[T EnvValueType](env *EnvData, name string, orElse T) T {
	result, err := GetEnv[T](env, name)
	if err != nil { return orElse }
	return result
}

func GetEnv[T EnvValueType](env *EnvData, name string) (T, error) {
	str, err := env.GetString(name)
	if err != nil {
		return GetDefault[T](), err
	}

	var zero T
	switch any(zero).(type) {
		case int8:
			num, err := strconv.ParseInt(str, 10, 8)
			if err != nil {
				return zero, err
			}
			return any(int8(num)).(T), nil
		case uint8:
			num, err := strconv.ParseUint(str, 10, 8)
			if err != nil {
				return zero, err
			}
			return any(uint8(num)).(T), nil
		case int16:
			num, err := strconv.ParseInt(str, 10, 16)
			if err != nil {
				return zero, err
			}
			return any(int16(num)).(T), nil
		case uint16:
			num, err := strconv.ParseUint(str, 10, 16)
			if err != nil {
				return zero, err
			}
			return any(uint16(num)).(T), nil
		case int32:
			num, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				return zero, err
			}
			return any(int32(num)).(T), nil
		case uint32:
			num, err := strconv.ParseUint(str, 10, 32)
			if err != nil {
				return zero, err
			}
			return any(uint32(num)).(T), nil
		case int64:
			num, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return zero, err
			}
			return any(int64(num)).(T), nil
		case uint64:
			num, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return zero, err
			}
			return any(uint64(num)).(T), nil
		case float32:
			num, err := strconv.ParseFloat(str, 32)
			if err != nil {
				return zero, err
			}
			return any(float32(num)).(T), nil
		case float64:
			num, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return zero, err
			}
			return any(float64(num)).(T), nil
		case bool:
			num, err := strconv.ParseBool(str)
			if err != nil {
				return zero, err
			}
			return any(bool(num)).(T), nil
		case string:
			return any(str).(T), nil
	}

	return zero, fmt.Errorf("unsupported type")
}
