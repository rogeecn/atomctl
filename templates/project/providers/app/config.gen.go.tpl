// Code generated by go-enum DO NOT EDIT.
// Version: -
// Revision: -
// Build Date: -
// Built By: -

package app

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// AppModeDevelopment is a AppMode of type development.
	AppModeDevelopment AppMode = "development"
	// AppModeRelease is a AppMode of type release.
	AppModeRelease AppMode = "release"
	// AppModeTest is a AppMode of type test.
	AppModeTest AppMode = "test"
)

var ErrInvalidAppMode = fmt.Errorf("not a valid AppMode, try [%s]", strings.Join(_AppModeNames, ", "))

var _AppModeNames = []string{
	string(AppModeDevelopment),
	string(AppModeRelease),
	string(AppModeTest),
}

// AppModeNames returns a list of possible string values of AppMode.
func AppModeNames() []string {
	tmp := make([]string, len(_AppModeNames))
	copy(tmp, _AppModeNames)
	return tmp
}

// AppModeValues returns a list of the values for AppMode
func AppModeValues() []AppMode {
	return []AppMode{
		AppModeDevelopment,
		AppModeRelease,
		AppModeTest,
	}
}

// String implements the Stringer interface.
func (x AppMode) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x AppMode) IsValid() bool {
	_, err := ParseAppMode(string(x))
	return err == nil
}

var _AppModeValue = map[string]AppMode{
	"development": AppModeDevelopment,
	"release":     AppModeRelease,
	"test":        AppModeTest,
}

// ParseAppMode attempts to convert a string to a AppMode.
func ParseAppMode(name string) (AppMode, error) {
	if x, ok := _AppModeValue[name]; ok {
		return x, nil
	}
	return AppMode(""), fmt.Errorf("%s is %w", name, ErrInvalidAppMode)
}

var errAppModeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *AppMode) Scan(value interface{}) (err error) {
	if value == nil {
		*x = AppMode("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseAppMode(v)
	case []byte:
		*x, err = ParseAppMode(string(v))
	case AppMode:
		*x = v
	case *AppMode:
		if v == nil {
			return errAppModeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errAppModeNilPtr
		}
		*x, err = ParseAppMode(*v)
	default:
		return errors.New("invalid type for AppMode")
	}

	return
}

// Value implements the driver Valuer interface.
func (x AppMode) Value() (driver.Value, error) {
	return x.String(), nil
}

// Set implements the Golang flag.Value interface func.
func (x *AppMode) Set(val string) error {
	v, err := ParseAppMode(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *AppMode) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *AppMode) Type() string {
	return "AppMode"
}

type NullAppMode struct {
	AppMode AppMode
	Valid   bool
}

func NewNullAppMode(val interface{}) (x NullAppMode) {
	err := x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	_ = err            // make any errcheck linters happy
	return
}

// Scan implements the Scanner interface.
func (x *NullAppMode) Scan(value interface{}) (err error) {
	if value == nil {
		x.AppMode, x.Valid = AppMode(""), false
		return
	}

	err = x.AppMode.Scan(value)
	x.Valid = (err == nil)
	return
}

// Value implements the driver Valuer interface.
func (x NullAppMode) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	// driver.Value accepts int64 for int values.
	return string(x.AppMode), nil
}

type NullAppModeStr struct {
	NullAppMode
}

func NewNullAppModeStr(val interface{}) (x NullAppModeStr) {
	x.Scan(val) // yes, we ignore this error, it will just be an invalid value.
	return
}

// Value implements the driver Valuer interface.
func (x NullAppModeStr) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	return x.AppMode.String(), nil
}
