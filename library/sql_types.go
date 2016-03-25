package library

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
)

// Language represents an ISO639 language
type Language struct {
	base language.Base
}

// NewLanguage creates a language from language.Base
func NewLanguage(base language.Base) Language {
	return Language{
		base: base,
	}
}

// ParseLanguage parses a language code (either 2-letter or 3-letter)
func ParseLanguage(code string) (Language, error) {
	l := Language{}
	var err error
	l.base, err = language.ParseBase(code)
	return l, err
}

// String returns the 2-letter representation of the language
// (same as ISO2, implements the fmt.Stringer interface)
func (l *Language) String() string {
	return l.ISO2()
}

// ISO2 returns the 2-letter representation of the language
func (l *Language) ISO2() string {
	return l.base.String()
}

// ISO3 returns the 3-letter representation of the language
func (l *Language) ISO3() string {
	return l.base.ISO3()
}

// Languages represents an array of languages
// (its SQL type should be text or varchar)
type Languages []Language

// NewLanguages creates a list of languages based on bases
func NewLanguages(bases []language.Base) Languages {
	languages := make(Languages, len(bases))
	for i := range bases {
		languages[i] = NewLanguage(bases[i])
	}
	return languages
}

// ParseLanguages parses languages from a space-separated list of language codes
func ParseLanguages(spaceSeparatedCodes string) (Languages, error) {
	codes := strings.Split(spaceSeparatedCodes, " ")
	languages := make(Languages, len(codes))

	var (
		err error
		j   int
	)

	for i := range codes {
		if codes[i] != "" {
			languages[j], err = ParseLanguage(codes[i])
			if err != nil {
				return nil, err
			}
			j++
		}
	}

	return languages, nil
}

// String represents the languages as 2-letter codes separated by spaces
func (l Languages) String() string {
	codes := make([]string, len(l))
	for i := range l {
		codes[i] = l[i].String()
	}
	return strings.Join(codes, " ")
}

// Scan deserialises the object from raw database data
func (l *Languages) Scan(src interface{}) error {
	var (
		err   error
		codes string
	)

	switch data := src.(type) {
	case string:
		codes = data
	case []byte:
		codes = string(data)
	default:
		return fmt.Errorf("unknown languages type")
	}

	*l, err = ParseLanguages(codes)
	if err != nil {
		return err
	}
	return nil
}

// Value serialises the object to raw database data
func (l Languages) Value() (driver.Value, error) {
	return l.String(), nil
}

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
