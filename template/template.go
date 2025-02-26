// Package template provides a simple wrapper around the text/template package
// for rendering templates with data. It supports loading templates from
// embedded files and provides methods for rendering templates to bytes,
// strings, and io.Reader.
package template

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"sync"
	"text/template"
)

//go:embed **/*.tmpl
var FS embed.FS

// Engine is a wrapper around the text/template.Template type. It provides
// methods for rendering templates with data.
type Engine struct {
	mu        sync.RWMutex
	templates map[string]*template.Template
}

// New creates a new template engine. It parses all templates in the templates
// directory and returns an Engine instance. If parsing fails, it returns
// an error.
func New() *Engine {
	return &Engine{
		templates: make(map[string]*template.Template),
	}
}

// lookup parses a template from the embedded file system. It caches the parsed
// template in the engine's templates map for future use. If parsing fails, it
// returns an error.
func (e *Engine) lookup(name string) (*template.Template, error) {
	tmpl, err := template.ParseFS(FS, name)
	if err != nil {
		e.templates[name] = nil
		return nil, fmt.Errorf("failed to parse template %q: %w", name, err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	return tmpl, nil
}

// Lookup returns a template by name. If the template is not found, it returns
// an error.
func (e *Engine) Lookup(name string) (*template.Template, error) {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if !ok {
		return e.lookup(name)
	}
	return tmpl, nil
}

// RenderBytes returns a byte slice for the rendered template.
func (e *Engine) RenderBytes(name string, data any) ([]byte, error) {
	tmpl, err := e.Lookup(name)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	if err = tmpl.Execute(&result, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return result.Bytes(), nil
}

// RenderString returns a string for the rendered template.
func (e *Engine) RenderString(name string, data any) (string, error) {
	b, err := e.RenderBytes(name, data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Render returns an io.Reader for the rendered template.
func (e *Engine) Render(name string, data any) (io.Reader, error) {
	b, err := e.RenderBytes(name, data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
