package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// BigUint64 is an encapsulated uint64 that can be stored in a database
type BigUint64 uint64

// Scan deserialises the object from raw database data
func (b *BigUint64) Scan(src interface{}) error {
	var (
		intText string
		err     error
	)

	switch data := src.(type) {
	case string:
		intText = data
	case []byte:
		intText = string(data)
	default:
		return fmt.Errorf("unknown uint64 type")
	}

	value, err := strconv.ParseUint(intText, 16, 64)
	if err != nil {
		return err
	}
	*b = BigUint64(value)
	return nil
}

// Value serialises the object to raw database data
func (b BigUint64) Value() (driver.Value, error) {
	return fmt.Sprintf("%016x", uint64(b)), nil
}

// Duration is an instance of time.Duration which implements the SQL
// Valuer and Scanner interfaces, so it can be stored in a database.
// Its SQL type should be integer.
type Duration time.Duration

// Scan deserialises the object from raw database data
func (d *Duration) Scan(src interface{}) error {
	switch data := src.(type) {
	case int64:
		*d = Duration(data)
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
// Its SQL type should be blob.
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

// MapStringStepStatus is an instance of map[string]*StepStatus which
// implements the SQL Valuer and Scanner interfaces, so it can be stored
// in a database. Its SQL type should be blob.
type MapStringStepStatus map[string]*StepStatus

// Scan deserialises the object from raw database data
func (m *MapStringStepStatus) Scan(src interface{}) error {
	switch data := src.(type) {
	case []byte:
		result := make(MapStringStepStatus)
		err := json.Unmarshal(data, &result)
		if err != nil {
			return fmt.Errorf("unable to parse map: %s", err)
		}
		*m = result
	default:
		return fmt.Errorf("unknown type for map[string]*StepStatus")
	}
	return nil
}

// Value serialises the object to raw database data
func (m MapStringStepStatus) Value() (driver.Value, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return nil, fmt.Errorf("unable to serialise map: %s", err)
	}
	return data, nil
}

// For returns the status for the given step (key). Same as using [], but
// when the key is not found it fills it with a default status instead
// of returning nil
func (m MapStringStepStatus) For(key string) *StepStatus {
	status := m[key]
	if status == nil {
		status = &StepStatus{}
		m[key] = status
	}
	return status
}

// Status represents the state of an import step
type Status int

//go:generate stringer -type=Status

const (
	// Incomplete objects haven't been imported yet
	Incomplete Status = iota
	// Skipped objects have been skipped by user action
	Skipped
	// Success objects have been imported correctly
	Success
	// Error objects have failed importing
	Error
)

// StepStatus contains the state of an import state and a string message
type StepStatus struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

// Errorf sets the status to Error and writes a message
func (m *StepStatus) Errorf(message string, arguments ...interface{}) {
	m.Status = Error
	m.Message = fmt.Sprintf(message, arguments...)
}

// Succeed sets the status to Success
func (m *StepStatus) Succeed() {
	m.Status = Success
	m.Message = ""
}

// Skip sets the status to Skipped
func (m *StepStatus) Skip() {
	m.Status = Skipped
	m.Message = ""
}

// Errorf creates a string pointer to an error message
func Errorf(message string, arguments ...interface{}) *string {
	errorMessage := fmt.Sprintf(message, arguments...)
	return &errorMessage
}
