package msg

import (
	"bytes"
	"html/template"
	"log/slog"
)

// port
type Renderer interface {
	Render(string, any) ([]byte, error)
}

// adapter
type RendererStdlib struct {
	Registry *template.Template
	Log      *slog.Logger
}

func (r *RendererStdlib) Render(name string, data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := r.Registry.ExecuteTemplate(buf, name, data)
	if err != nil {
		r.Log.Error("rendering failed", slog.Any("reason", err))
	}
	return buf.Bytes(), err
}
