package apidoc

import (
	"fmt"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	apiAttr         = "@api"
	titleAttr       = "@title"
	versionAttr     = "@version"
	descriptionAttr = "@desc"
	successAttr     = "@success"
	failureAttr     = "@failure"
	responseAttr    = "@response"
	deprecatedAttr  = "@deprecated"
	tagsAttr        = "@tags"
	authorAttr      = "@author"

	//doc
	hostAttr     = "@host"
	basePathAttr = "@basepath"
)

var allMethod = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPut:     {},
	http.MethodPost:    {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
	http.MethodHead:    {},
	http.MethodPatch:   {},
}

type Parser struct {
	doc      *ApiDocSpec
	packages *PackagesDefinitions
	// excludes excludes dirs and files in SearchDir
	excludes map[string]struct{}
}

func New() *Parser {
	return &Parser{
		doc:      &ApiDocSpec{},
		packages: NewPackagesDefinitions(),
		excludes: make(map[string]struct{}),
	}
}

func SetExcludedDirsAndFiles(excludes string) func(*Parser) {
	return func(p *Parser) {
		for _, f := range strings.Split(excludes, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				f = filepath.Clean(f)
				p.excludes[f] = struct{}{}
			}
		}
	}
}

func (parser *Parser) Parse(searchDir string, mainFile string) error {
	packageDir, err := getPkgName(searchDir)
	if err != nil {
		return err
	}
	if err = parser.getAllGoFileInfo(packageDir, searchDir); err != nil {
		return err
	}
	mainPath, err := filepath.Abs(filepath.Join(searchDir, mainFile))
	if err != nil {
		return err
	}
	if err = parser.parseApiDocInfo(mainPath); err != nil {
		return err
	}
	if err = rangeFiles(parser.packages.files, parser.parseApiInfos); err != nil {
		return err
	}
	return nil
}

func (parser *Parser) GetApiDoc() *ApiDocSpec {
	return parser.doc
}

func (parser *Parser) parseApiDocInfo(mainPath string) error {
	fileTree, err := goparser.ParseFile(token.NewFileSet(), mainPath, nil, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse source files %s: %s", mainPath, err)
	}
	for _, comment := range fileTree.Comments {
		comments := strings.Split(comment.Text(), "\n")
		if !isApiDocComment(comments) {
			continue
		}

		err = parseApiDocInfo(parser, comments)
		if err != nil {
			return err
		}
	}

	return nil
}

func (parser *Parser) parseApiInfos(fileName string, astFile *ast.File) error {
	for _, astDescription := range astFile.Decls {
		astDeclaration, ok := astDescription.(*ast.FuncDecl)
		if ok && astDeclaration.Doc != nil && astDeclaration.Doc.List != nil {
			if astDeclaration.Name.Name == "main" {
				continue
			}
			operation := NewOperation(parser)
			for _, comment := range astDeclaration.Doc.List {
				err := operation.ParseComment(comment.Text, astFile)
				if err != nil {
					return fmt.Errorf("ParseComment error in file %s :%+v", fileName, err)
				}
			}
			parser.doc.Apis = append(parser.doc.Apis, operation.ApiSpec)
		}
	}

	return nil
}

func parseApiDocInfo(parser *Parser, comments []string) error {
	previousAttribute := ""
	for line := 0; line < len(comments); line++ {
		commentLine := comments[line]
		attribute := strings.Split(commentLine, " ")[0]
		value := strings.TrimSpace(commentLine[len(attribute):])
		multilineBlock := false
		if previousAttribute == attribute {
			multilineBlock = true
		}
		switch strings.ToLower(attribute) {
		case versionAttr:
			parser.doc.Version = value
		case titleAttr:
			parser.doc.Title = value
		case descriptionAttr:
			if multilineBlock {
				parser.doc.Description += "\n" + value
				continue
			}
			parser.doc.Description = value
		case basePathAttr:
			parser.doc.BasePath = value
		case hostAttr:
			parser.doc.Host = value
		}
	}
	return nil
}

func isApiDocComment(comments []string) bool {
	for _, commentLine := range comments {
		attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
		switch attribute {
		case apiAttr, successAttr, failureAttr, responseAttr:
			return false
		}
	}

	return true
}

func getPkgName(searchDir string) (string, error) {
	cmd := exec.Command("go", "list", "-f={{.ImportPath}}")
	cmd.Dir = searchDir

	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list command, %s, stdout:%s, stderr:%s", err, stdout.String(), stderr.String())
	}

	outStr, _ := stdout.String(), stderr.String()

	if outStr[0] == '_' { // will shown like _/{GOPATH}/src/{YOUR_PACKAGE} when NOT enable GO MODULE.
		outStr = strings.TrimPrefix(outStr, "_"+build.Default.GOPATH+"/src/")
	}

	f := strings.Split(outStr, "\n")

	outStr = f[0]

	return outStr, nil
}

// GetAllGoFileInfo gets all Go source files information for given searchDir.
func (parser *Parser) getAllGoFileInfo(packageDir, searchDir string) error {
	return filepath.Walk(searchDir, func(path string, f os.FileInfo, _ error) error {
		err := parser.Skip(path, f)
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(searchDir, path)
		if err != nil {
			return err
		}

		return parser.parseFile(filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(packageDir, relPath)))), path, nil)
	})
}

func (parser *Parser) parseFile(packageDir, path string, src interface{}) error {
	if strings.HasSuffix(strings.ToLower(path), "_test.go") || filepath.Ext(path) != ".go" {
		return nil
	}

	// positions are relative to FileSet
	astFile, err := goparser.ParseFile(token.NewFileSet(), path, src, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("ParseFile error:%+v", err)
	}

	err = parser.packages.CollectAstFile(packageDir, path, astFile)
	if err != nil {
		return err
	}

	return nil
}

// Skip returns filepath.SkipDir error if match vendor and hidden folder.
func (parser *Parser) Skip(path string, f os.FileInfo) error {
	return walkWith(parser.excludes)(path, f)
}

func walkWith(excludes map[string]struct{}) func(path string, fileInfo os.FileInfo) error {
	return func(path string, f os.FileInfo) error {
		if f.IsDir() {
			if f.Name() == "vendor" || // ignore "vendor"
				len(f.Name()) > 1 && f.Name()[0] == '.' { // exclude all hidden folder
				return filepath.SkipDir
			}

			if excludes != nil {
				if _, ok := excludes[path]; ok {
					return filepath.SkipDir
				}
			}
		}

		return nil
	}
}

func fullTypeName(pkgName, typeName string) string {
	if pkgName != "" {
		return pkgName + "." + typeName
	}

	return typeName
}
