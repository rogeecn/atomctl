package fields

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/samber/lo"
)

// implement sql.Scanner interface
type field struct{}

func (x *field) Scan(value interface{}) (err error) {
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

func (x field) Value() (driver.Value, error) {
	return json.Marshal(x)
}

func (x field) MustValue() driver.Value {
	return lo.Must(json.Marshal(x))
}
