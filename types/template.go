package types

import (
	"bytes"
	"text/template"
)

// Template wraps a text template for generating strings
type Template struct {
	template.Template
}

// UnmarshalText parses the template from text
func (t *Template) UnmarshalText(text []byte) error {
	templ := template.New("template")
	templ, err := templ.Parse(string(text))
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
