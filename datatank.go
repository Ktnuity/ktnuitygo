package ktnuitygo

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveJson[T any](filename string, data T) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

type TankLoadError struct {
	isSafe			bool
	message			string
}

func (err *TankLoadError) Error() string {
	return err.message
}

func loadJson[T any](filename string, target *T) *TankLoadError {
	file, err := os.Open(filename)
	if err != nil {
		return &TankLoadError{
			isSafe: true,
			message: fmt.Sprintf("failed to open file '%s': %v", filename, err),
		}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(target); err != nil {
		return &TankLoadError{
			isSafe: false,
			message: fmt.Sprintf("failed to decode JSON: %v", err),
		}
	}

	return nil
}

var tankDir string = "."
func DataTankSetDir(dir string) {
	for len(dir) > 0 && dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}

	if len(dir) == 0 {
		tankDir = "."
		return
	}

	tankDir = dir
}

func DataTankNew[T any](name string) (*DataTank[T], error) {
	data := InitDefault[T]()
	
	err := loadJson(tankPath(name), &data)
	if err != nil {
		if !err.isSafe {
			return nil, fmt.Errorf("failed to load DataTank '%s' data: %w", name, err)
		}

		data = InitDefault[T]()
	}

	return &DataTank[T]{
		name: name,
		data: &data,
	}, nil
}

type DataTankSetFn[T any] func(data *T)
type DataTankGetFn[R any, T any] func(data *T) *R

type DataTank[T any] struct {
	name string
	data *T
}

func tankPath(name string) string {
	return fmt.Sprintf("%s/%s.tank.json", tankDir, name)
}

func (d *DataTank[T]) Save() error {
	return saveJson(tankPath(d.name), d.data)
}

func (d *DataTank[T]) Reload() error {
	data := GetDefault[T]()
	err := loadJson(tankPath(d.name), &data)
	if err != nil {
		return fmt.Errorf("failed to reload DataTank '%s' data: %w", d.name, err)
	}

	d.data = &data
	return nil
}

func DataTankGet[R any, T any](d *DataTank[T], fn DataTankGetFn[R, T]) *R {
	return fn(d.data)
}

func DataTankSet[T any](d *DataTank[T], fn DataTankSetFn[T]) error {
	fn(d.data)
	return d.Save()
}


