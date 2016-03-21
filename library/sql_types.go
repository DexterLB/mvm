package library

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Duration is an instance of time.Duration which implements the SQL
// Valuer and Scanner interfaces, so it can be stored in a database.
// Its SQL type should be integer.
type Duration time.Duration

// Scan deserialises the object from raw database data
func (d *Duration) Scan(src interface{}) error {
	switch data := src.(type) {
	case int64:
		*d = Duration(Duration(data))
	default:
		return fmt.Errorf("unknown duration type")
	}
	return nil
}

// Value serialises the object to raw database data
func (d Duration) Value() (driver.Value, error) {
	return int64(d), nil
}

// MapStringString is an instance of map[string]string which implements the SQL
// Valuer and Scanner interfaces, so it can be stored in a database.
// Its SQL type should be integer.
type MapStringString map[string]string

// Scan deserialises the object from raw database data
func (m *MapStringString) Scan(src interface{}) error {
	switch data := src.(type) {
	case []byte:
		result := make(map[string]string)
		err := json.Unmarshal(data, &result)
		if err != nil {
			return fmt.Errorf("unable to parse map: %s", err)
		}
		*m = result
	default:
		return fmt.Errorf("unknown type for map[string]string")
	}
	return nil
}

// Value serialises the object to raw database data
func (m MapStringString) Value() (driver.Value, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return nil, fmt.Errorf("unable to serialise map: %s", err)
	}
	return data, nil
}

// SliceString is an instance of []string which implements the SQL
// Valuer and Scanner interfaces, so it can be stored in a database.
// Its SQL type should be integer.
type SliceString []string

// Scan deserialises the object from raw database data
func (m *SliceString) Scan(src interface{}) error {
	switch data := src.(type) {
	case []byte:
		var result []string
		err := json.Unmarshal(data, &result)
		if err != nil {
			return fmt.Errorf("unable to parse slice: %s", err)
		}
		*m = result
	default:
		return fmt.Errorf("unknown type for []string")
	}
	return nil
}

// Value serialises the object to raw database data
func (m SliceString) Value() (driver.Value, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return nil, fmt.Errorf("unable to serialise slice: %s", err)
	}
	return data, nil
}
