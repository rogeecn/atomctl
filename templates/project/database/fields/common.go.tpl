package fields

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// implement sql.Scanner interface
type Json[T any] struct {
	Data T `json:",inline"`
}

func ToJson[T any](data T) Json[T] {
	return Json[T]{Data: data}
}

func (x *Json[T]) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), &x)
	case []byte:
		return json.Unmarshal(v, &x)
	case *string:
		return json.Unmarshal([]byte(*v), &x)
	}
	return errors.New("Unknown type for ")
}

func (x Json[T]) Value() (driver.Value, error) {
	return json.Marshal(x.Data)
}

func (x Json[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Data)
}

func (x *Json[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	x.Data = value
	return nil
}
