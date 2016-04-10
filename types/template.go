package types

import (
	"bytes"
	"text/template"
)

// Template wraps a text template for generating strings
type Template struct {
	template.Template
}

// ParseTemplate creates a new template from the string
func ParseTemplate(s string) (*Template, error) {
	t := &Template{}
	err := t.UnmarshalString(s)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ParseTemplate creates a new template from the string and panics on failiure.
func MustParseTemplate(s string) *Template {
	t, err := ParseTemplate(s)
	if err != nil {
		panic(err)
	}
	return t
}

// UnmarshalText parses the template from text
func (t *Template) UnmarshalText(text []byte) error {
	return t.UnmarshalString(string(text))
}

// UnmarshalString parses the template from a string
func (t *Template) UnmarshalString(s string) error {
	templ := template.New("template")
	templ, err := templ.Parse(s)
	if err != nil {
		return err
	}
	t.Template = *templ
	return nil
}

// On performs the template on data
func (t *Template) On(data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
