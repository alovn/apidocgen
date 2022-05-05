package gen

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/alovn/apidoc"
)

type Gen struct {
}

func New() *Gen {
	return &Gen{}
}

type Config struct {
	SearchDir    string
	OutputDir    string
	TemplateFile string
	MainFile     string
}

func (g *Gen) Build(c *Config) error {
	if c == nil {
		return errors.New("error config")
	}
	p := apidoc.New()
	if err := p.Parse(c.SearchDir, c.MainFile); err != nil {
		return err
	}
	doc := p.GetApiDoc()
	if err := os.MkdirAll(c.OutputDir, os.ModePerm); err != nil {
		return err
	}
	_ = doc

	t := template.New("docs")
	t, err := t.Parse(defaultTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(c.OutputDir, "api-docs.md"))
	if err != nil {
		return err
	}
	defer f.Close()
	if err = t.Execute(f, doc); err != nil {
		return err
	}
	return nil
}
