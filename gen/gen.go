package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/alovn/apidocgen/parser"
)

type Gen struct {
	c          *Config
	templateFS fs.FS
}

type Config struct {
	SearchDir         string
	OutputDir         string
	TemplateName      string
	CustomTemplateDir string
	ExcludesDir       string
	IsGenSingleFile   bool
}

type TemplateConfig struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

func New(c *Config) *Gen {
	var defaultTemplateName = "markdown"
	if c.TemplateName == "" {
		c.TemplateName = defaultTemplateName
	}
	if strings.ContainsAny(c.TemplateName, "/\\") {
		c.CustomTemplateDir = c.TemplateName
		c.TemplateName = ""
	}
	return &Gen{
		c: c,
	}
}

func (g *Gen) readTemplate(name string) (s string, err error) {
	var bs []byte
	if g.templateFS == nil {
		err = errors.New("error: templateFS nil")
		return
	}
	if bs, err = fs.ReadFile(g.templateFS, name); err != nil {
		return
	} else {
		s = string(bs)
		return
	}
}

func (g *Gen) Build() error {
	if g.c == nil {
		return errors.New("error config")
	}
	searchDirs := strings.Split(g.c.SearchDir, ",")
	for _, searchDir := range searchDirs {
		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
			return fmt.Errorf("dir: %s does not exist", searchDir)
		}
	}
	var err error
	if g.c.CustomTemplateDir != "" {
		g.templateFS = os.DirFS(g.c.CustomTemplateDir)
	} else {
		if g.templateFS, err = fs.Sub(defaultTemplateFS, fmt.Sprintf("template/%s", g.c.TemplateName)); err != nil {
			return err
		}
	}

	var templateSingleIndex string
	var templateGroupIndex string
	var templateGroupApis string
	var templateConfig TemplateConfig

	if s, err := g.readTemplate("config.json"); err != nil {
		return err
	} else {
		if err = json.Unmarshal([]byte(s), &templateConfig); err != nil {
			return err
		}
		if templateConfig.Index == "" {
			templateConfig.Index = "README.md"
		}
		fmt.Println("use template:", templateConfig.Name)
	}

	if g.c.IsGenSingleFile {
		if templateSingleIndex, err = g.readTemplate("single_index.tpl"); err != nil {
			return err
		}
	} else {
		if templateGroupIndex, err = g.readTemplate("group_index.tpl"); err != nil {
			return err
		}
		if templateGroupApis, err = g.readTemplate("group_apis.tpl"); err != nil {
			return err
		}
	}

	p := parser.New()
	parser.SetExcludedDirsAndFiles(g.c.ExcludesDir)(p)
	if err := p.Parse(searchDirs); err != nil {
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
		doc.Groups = append(doc.Groups, &parser.ApiGroupSpec{
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
	sortApis := func(apis []*parser.ApiSpec) {
		less := func(i, j int) bool {
			a, b := apis[i], apis[j]
			return a.Order < b.Order
		}
		sort.Slice(apis, less)
	}
	if g.c.IsGenSingleFile {
		t := template.New("single-index").Funcs(funcMap)
		t, err := t.Parse(templateSingleIndex)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(g.c.OutputDir, templateConfig.Index))
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
		fmt.Println("generated: ", templateConfig.Index)
	} else {
		//group
		t := template.New("group-apis").Funcs(funcMap)
		t, err := t.Parse(templateGroupApis)
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
		t, err = t.Parse(templateGroupIndex)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(g.c.OutputDir, templateConfig.Index))
		if err != nil {
			return err
		}
		defer f.Close()
		if err = t.Execute(f, doc); err != nil {
			return err
		}
		fmt.Println("generated:", templateConfig.Index)
	}
	fmt.Println("apis total count:", doc.TotalCount)
	return nil
}
