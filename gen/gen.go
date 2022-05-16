package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	SearchDir       string
	OutputDir       string
	TemplateFile    string
	IsGenSingleFile bool
}

func (g *Gen) Build() error {
	if g.c == nil {
		return errors.New("error config")
	}
	p := apidoc.New()
	if err := p.Parse(g.c.SearchDir); err != nil {
		return err
	}
	doc := p.GetApiDoc()
	if doc.Service == "" {
		fmt.Println("apidoc @service is not set")
		return nil
	}
	if doc.Title == "" {
		fmt.Println("apidoc @title is not set")
		return nil
	}
	if doc.TotalCount == 0 {
		fmt.Println("apis count is 0")
		return nil
	}

	if len(doc.UngroupedApis) > 0 {
		doc.Groups = append(doc.Groups, &apidoc.ApiGroupSpec{
			Group:       "ungrouped",
			Title:       "ungrouped",
			Description: "Ungrouped apis",
			Apis:        doc.UngroupedApis,
		})
		doc.UngroupedApis = doc.UngroupedApis[:0]
	}
	sort.Slice(doc.Groups, func(i, j int) bool {
		a, b := doc.Groups[i], doc.Groups[j]
		if a.Order == b.Order {
			return strings.Compare(a.Group, b.Group) < 0
		}
		return a.Order < b.Order
	})

	if err := os.MkdirAll(g.c.OutputDir, os.ModePerm); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
	sortApis := func(apis []*apidoc.ApiSpec) {
		less := func(i, j int) bool {
			a, b := apis[i], apis[j]
			return a.Order < b.Order
		}
		sort.Slice(apis, less)
	}

	if g.c.IsGenSingleFile {
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
			apis := g.Apis
			sortApis(apis)
		}

		if err = t.Execute(f, doc); err != nil {
			return err
		}
		fmt.Println("generated: README.md")
	} else {
		//group
		t := template.New("group-apis").Funcs(funcMap)
		t, err := t.Parse(groupApisTemplate)
		if err != nil {
			return err
		}
		for _, v := range doc.Groups {
			group := v
			sortApis(group.Apis)
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
		t = template.New("group-readme").Funcs(funcMap)
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
	}
	fmt.Println("apis total count:", doc.TotalCount)
	return nil
}
