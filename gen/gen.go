package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/alovn/apidoc"
)

type Gen struct {
	c *Config
}

func New(c *Config) *Gen {
	return &Gen{
		c: c,
	}
}

type Config struct {
	SearchDir      string
	OutputDir      string
	TemplateFile   string
	MainFile       string
	IsGenGroupFile bool
}

func (g *Gen) Build() error {
	if g.c == nil {
		return errors.New("error config")
	}
	p := apidoc.New()
	if err := p.Parse(g.c.SearchDir, g.c.MainFile); err != nil {
		return err
	}
	doc := p.GetApiDoc()
	if doc.Title == "" && doc.Service == "" && len(doc.Apis) == 0 {
		fmt.Println("can't find apis")
		return nil
	}

	if len(doc.Apis) > 0 {
		doc.Groups = append(doc.Groups, &apidoc.ApiGroupSpec{
			Group:       "ungrouped",
			Title:       "ungrouped",
			Description: "Ungrouped apis",
			Apis:        doc.Apis,
		})
	}

	if err := os.MkdirAll(g.c.OutputDir, os.ModePerm); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	if g.c.IsGenGroupFile {
		//group
		t := template.New("group").Funcs(funcMap)
		t, err := t.Parse(groupApisTemplate)
		if err != nil {
			return err
		}
		for _, v := range doc.Groups {
			group := v
			fileName := fmt.Sprintf("apis-%s.md", group.Group)
			f, err := os.Create(filepath.Join(g.c.OutputDir, fileName))
			if err != nil {
				return err
			}
			defer f.Close()
			if err = t.Execute(f, group); err != nil {
				return err
			}
			fmt.Println("Generated:", fileName)
		}

		//readme
		t = template.New("apis").Funcs(funcMap)
		t, err = t.Parse(groupReadmeTemplate)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(g.c.OutputDir, "README.md"))
		if err != nil {
			return err
		}
		defer f.Close()
		if err = t.Execute(f, doc); err != nil {
			return err
		}
		fmt.Println("generated: README.md")

		return nil
	}

	t := template.New("apis-single").Funcs(funcMap)
	t, err := t.Parse(singleApisTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(g.c.OutputDir, "README.md"))
	if err != nil {
		return err
	}
	defer f.Close()
	for _, g := range doc.Groups {
		doc.Apis = append(g.Apis, doc.Apis...)
	}
	if err = t.Execute(f, doc); err != nil {
		return err
	}
	fmt.Println("generated: README.md")
	return nil
}
