package generate

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/alexeyco/helmfile/assets"
	"github.com/alexeyco/helmfile/internal"
	"github.com/spf13/afero"
)

type Generator struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Generator {
	return &Generator{fs: fs}
}

func (g *Generator) Render(v internal.Versions) (string, error) {
	tmpl, err := template.New("Dockerfile").Parse(assets.DockerfileTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, v); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) Generate(v internal.Versions, path string) error {
	content, err := g.Render(v)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := afero.WriteFile(g.fs, path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
