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

	"github.com/alovn/apidoc"
)

type Gen struct {
	c          *Config
	templateFS fs.FS
	defaultFS  fs.FS
}

func New(c *Config) *Gen {
	return &Gen{
		c: c,
	}
}

type Config struct {
	SearchDir       string
	OutputDir       string
	TemplateDir     string
	Excludes        string
	IsGenSingleFile bool
}

type TemplateConfig struct {
	Name  string `json:"name"`
	Index string `json:"index"`
}

func (g *Gen) readTemplate(name string) (s string, err error) {
	var bs []byte
	if g.templateFS != nil {
		if bs, err = fs.ReadFile(g.templateFS, name); err != nil {
			return
		} else {
			s = string(bs)
			return
		}
	} else {
		if bs, err = fs.ReadFile(g.defaultFS, name); err != nil {
			return
		} else {
			s = string(bs)
			return
		}
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
	if g.defaultFS == nil {
		if g.defaultFS, err = fs.Sub(defaultTemplateFS, "template"); err != nil {
			return err
		}
	}

	if g.templateFS == nil && g.c.TemplateDir != "" {
		g.templateFS = os.DirFS(g.c.TemplateDir)
	}
	var templateSingle string
	var templateGroupIndex string
	var templateGroupApis string
	var templateConfig TemplateConfig

	if templateSingle, err = g.readTemplate("single.tpl"); err != nil {
		return err
	}
	if templateGroupIndex, err = g.readTemplate("group_index.tpl"); err != nil {
		return err
	}
	if templateGroupApis, err = g.readTemplate("group_apis.tpl"); err != nil {
		return err
	}
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

	p := apidoc.New()
	apidoc.SetExcludedDirsAndFiles(g.c.Excludes)(p)
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
		t, err := t.Parse(templateSingle)
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
