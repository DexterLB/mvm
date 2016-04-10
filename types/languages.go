package types

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// Language represents an ISO639 language. Its SQL type could be varchar(3).
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

// UnmarshalText parses a language code from raw text
func (l *Language) UnmarshalText(text []byte) error {
	lang, err := ParseLanguage(string(text))
	if err != nil {
		return err
	}
	*l = lang
	return nil
}

// MustParseLanguage parses a language code (either 2-letter or 3-letter).
// Panics on failiure.
func MustParseLanguage(code string) Language {
	language, err := ParseLanguage(code)
	if err != nil {
		panic(err)
	}
	return language
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

// MustParseLanguages parses languages from a space-separated list of
// language codes, and panics on failiure
func MustParseLanguages(spaceSeparatedCodes string) Languages {
	languages, err := ParseLanguages(spaceSeparatedCodes)
	if err != nil {
		panic(err)
	}
	return languages
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

// Scan deserialises the object from raw database data
func (l *Language) Scan(src interface{}) error {
	var (
		err  error
		code string
	)

	switch data := src.(type) {
	case string:
		code = data
	case []byte:
		code = string(data)
	default:
		return fmt.Errorf("unknown language type")
	}

	*l, err = ParseLanguage(code)
	if err != nil {
		return err
	}
	return nil
}

// Value serialises the object to raw database data
func (l Language) Value() (driver.Value, error) {
	return l.String(), nil
}

// Value serialises the object to raw database data
func (l Languages) Value() (driver.Value, error) {
	return l.String(), nil
}
