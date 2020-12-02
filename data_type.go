package gorm

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// Array returns the optimal driver.Valuer and sql.Scanner for an array or
// slice of any dimension.
func Array(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	return &Generic{a}
}

func Any(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	return &Generic{a}
}

// Generic implements the driver.Valuer and sql.Scanner interfaces for
// an array or slice of any dimension.
type Generic struct{ A interface{} }

// Scan implements the sql.Scanner interface.
func (a *Generic) Scan(src interface{}) error {
	if a == nil {
		return fmt.Errorf("GenericStruct.Scan: %s", "a is nil")
	}
  if src == nil {
    return nil
  }
	dpv := reflect.ValueOf(a.A)
	switch {
	case dpv.Kind() != reflect.Ptr:
		return fmt.Errorf("pq: destination %T is not a pointer to array or slice", a.A)
	case dpv.IsNil():
		return fmt.Errorf("pq: destination %T is nil", a.A)
	}

	return json.Unmarshal(reflect.ValueOf(src).Bytes(), a.A)
}

// Value implements the driver.Valuer interface.
func (a Generic) Value() (driver.Value, error) {
	if a.A == nil {
		return nil, nil
	}
	return json.Marshal(a.A)
}
