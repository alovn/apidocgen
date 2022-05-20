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
	SearchDir       string
	OutputDir       string
	OutputIndexName string
	TemplateDir     string
	ExcludesDir     string
	IsGenSingleFile bool
}

type TemplateConfig struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

func New(c *Config) *Gen {
	var defaultTemplateName = "markdown"
	if c.TemplateDir == "" {
		c.TemplateDir = defaultTemplateName
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
		return errors.New("error: config nil")
	}

	searchDirs := strings.Split(g.c.SearchDir, ",")
	for _, searchDir := range searchDirs {
		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
			return fmt.Errorf("error: dir %s does not exist", searchDir)
		}
	}

	var err error
	if strings.ContainsAny(g.c.TemplateDir, "/\\") { //custom dir
		g.templateFS = os.DirFS(g.c.TemplateDir)
	} else {
		g.c.TemplateDir = fmt.Sprintf("template/%s", g.c.TemplateDir)
		if g.templateFS, err = fs.Sub(defaultTemplateFS, g.c.TemplateDir); err != nil {
			return err
		}
	}

	if s, err := g.readTemplate("config.json"); err == nil {
		var templateConfig TemplateConfig
		if err = json.Unmarshal([]byte(s), &templateConfig); err != nil {
			return err
		}
		if templateConfig.Index == "" {
			templateConfig.Index = "README.md"
		}
		if g.c.OutputIndexName == "" {
			g.c.OutputIndexName = templateConfig.Index
		}
		fmt.Println("use template:", templateConfig.Name)
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

	g.c.OutputIndexName = strings.ReplaceAll(g.c.OutputIndexName, "@{service}", doc.Service)
	if strings.ContainsAny(g.c.OutputIndexName, "/\\") { //custom dir
		return fmt.Errorf("error: output-index can't be a directory: %s.", g.c.OutputIndexName)
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

	t, err := template.New("apidocgen").Funcs(funcMap).ParseFS(g.templateFS, "*.tpl")
	if err != nil {
		return err
	}
	if g.c.IsGenSingleFile {
		f, err := os.Create(filepath.Join(g.c.OutputDir, g.c.OutputIndexName))
		if err != nil {
			return err
		}
		defer f.Close()
		for _, g := range doc.Groups {
			apis := g.Apis
			sortApis(apis)
		}
		if err = t.ExecuteTemplate(f, "single_index", doc); err != nil {
			return err
		}
		fmt.Println("generated: ", g.c.OutputIndexName)
	} else {
		//group
		for _, v := range doc.Groups {
			group := v
			sortApis(group.Apis)
			fileName := fmt.Sprintf("apis-%s.md", group.Group)
			f, err := os.Create(filepath.Join(g.c.OutputDir, fileName))
			if err != nil {
				return err
			}
			defer f.Close()
			if err = t.ExecuteTemplate(f, "group_apis", group); err != nil {
				return err
			}
			fmt.Println("Generated:", fileName)
		}
		//readme
		f, err := os.Create(filepath.Join(g.c.OutputDir, g.c.OutputIndexName))
		if err != nil {
			return err
		}
		defer f.Close()
		if err = t.ExecuteTemplate(f, "group_index", doc); err != nil {
			return err
		}
		fmt.Println("generated:", g.c.OutputIndexName)
	}
	fmt.Println("apis total count:", doc.TotalCount)
	return nil
}
