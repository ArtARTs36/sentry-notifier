package template

import (
	"bytes"

	"github.com/tyler-sommer/stick"
)

type Renderer struct {
	loader *stick.MemoryLoader
	engine *stick.Env
}

func NewRenderer(templates map[string]string) *Renderer {
	loader := &stick.MemoryLoader{
		Templates: templates,
	}

	return &Renderer{
		loader: loader,
		engine: stick.New(loader),
	}
}

func (r *Renderer) Render(templateKey string, params map[string]stick.Value) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := r.engine.Execute(templateKey, buf, params)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), err
}
