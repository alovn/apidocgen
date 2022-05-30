package gen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/alovn/apidocgen/mock"
	"github.com/alovn/apidocgen/parser"
)

type Gen struct {
	c          *Config
	templateFS fs.FS
	searchDirs []string
	parser     *parser.Parser
	mockApis   map[string][]mock.MockAPI
}

type Config struct {
	SearchDir        string
	OutputDir        string
	OutputIndexName  string
	TemplateDir      string
	ExcludesDir      string
	IsGenSingleFile  bool
	IsGenMocks       bool
	IsMockServer     bool
	MockServerListen string
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
		c:        c,
		mockApis: map[string][]mock.MockAPI{},
	}
}

func (g *Gen) Build() error {
	if g.c == nil {
		return errors.New("error: config nil")
	}
	g.searchDirs = strings.Split(g.c.SearchDir, ",")
	for _, searchDir := range g.searchDirs {
		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
			return fmt.Errorf("error: dir %s does not exist", searchDir)
		}
	}
	//parser
	g.parser = parser.New()
	parser.SetExcludedDirsAndFiles(g.c.ExcludesDir)(g.parser)
	if err := g.parser.Parse(g.searchDirs); err != nil {
		return err
	}
	doc := g.parser.GetApiDoc()
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
	if err := g.buildDocs(); err != nil {
		return err
	}
	if g.c.IsGenMocks || g.c.IsMockServer {
		if err := g.buildMocks(); err != nil {
			return err
		}
	}
	if g.c.IsGenMocks {
		if err := g.genMocks(); err != nil {
			return err
		}
	}
	if g.c.IsMockServer {
		mockServer := mock.New(g.c.MockServerListen)
		for _, mockApis := range g.mockApis {
			mockServer.InitMockApis(mockApis)
		}
		return mockServer.Serve()
	}
	return nil
}

func (g *Gen) buildDocs() error {
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

	doc := g.parser.GetApiDoc()

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
		defer f.Close() //#nosec
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
			f, err := os.Create(filepath.Clean(filepath.Join(g.c.OutputDir, fileName)))
			if err != nil {
				return err
			}
			defer f.Close() //#nosec
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
		defer f.Close() //#nosec
		if err = t.ExecuteTemplate(f, "group_index", doc); err != nil {
			return err
		}
		fmt.Println("generated:", g.c.OutputIndexName)
	}
	fmt.Println("apis total count:", doc.TotalCount)
	return nil
}

func (g *Gen) buildMocks() error {
	doc := g.parser.GetApiDoc()
	var basePath string
	if doc.BaseURL != "" {
		u, err := url.Parse(doc.BaseURL)
		if err != nil {
			return err
		}
		basePath = u.Path
	}

	for _, group := range doc.Groups {
		var mockApis []mock.MockAPI
		for _, api := range group.Apis {
			var contentType string
			switch strings.ToLower(api.Format) {
			case "xml":
				contentType = "application/xml"
			default:
				contentType = "application/json"
			}

			mapi := mock.MockAPI{
				Title:      api.Title,
				HTTPMethod: api.HTTPMethod,
				Path:       fmt.Sprintf("%s/%s", strings.TrimSuffix(basePath, "/"), strings.TrimPrefix(api.Api, "/")),
				Headers: map[string]string{
					"Content-Type": contentType,
				},
			}

			for _, res := range api.Responses {
				mapi.Responses = append(mapi.Responses, mock.MockAPIResponse{
					HTTPCode: res.StatusCode,
					Body:     res.PureBody(),
					IsMock:   res.IsMock || len(api.Responses) == 1,
				})
			}
			mockApis = append(mockApis, mapi)
		}
		g.mockApis[group.Group] = mockApis
	}
	return nil
}

func (g *Gen) genMocks() error {
	mocksDir := filepath.Join(g.c.OutputDir, "mocks")
	if err := os.MkdirAll(mocksDir, os.ModePerm); err != nil {
		return err
	}
	genMockFile := func(name string, bytes []byte) error {
		mocksFileName := fmt.Sprintf("%s.mocks", strings.ToLower(name))
		f, err := os.Create(filepath.Clean(filepath.Join(mocksDir, mocksFileName)))
		if err != nil {
			return err
		}
		defer f.Close() //#nosec
		if _, err := f.Write(bytes); err != nil {
			return err
		}
		return nil
	}

	for groupName, mockApis := range g.mockApis {
		if bytes, err := json.MarshalIndent(mockApis, "", " "); err == nil {
			if err := genMockFile(groupName, bytes); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
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
