package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	serviceAttr     = "@service"
	apiAttr         = "@api"
	titleAttr       = "@title"
	groupAttr       = "@group"
	versionAttr     = "@version"
	descriptionAttr = "@desc"
	acceptAttr      = "@accept"
	requestAttr     = "@request"
	queryAttr       = "@query"
	paramAttr       = "@param"
	headerAttr      = "@header"
	formAttr        = "@form"
	successAttr     = "@success"
	failureAttr     = "@failure"
	responseAttr    = "@response"
	formatAttr      = "@format"
	deprecatedAttr  = "@deprecated"
	authorAttr      = "@author"
	orderAttr       = "@order" //for sort

	//doc
	baseURLAttr = "@baseurl"
)

var allMethod = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPut:     {},
	http.MethodPost:    {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
	http.MethodHead:    {},
	http.MethodPatch:   {},
	"ANY":              {},
}

type Parser struct {
	doc      *ApiDocSpec
	groups   map[string]*ApiGroupSpec
	packages *PackagesDefinitions
	// excludes excludes dirs and files in SearchDir
	excludes map[string]struct{}
}

func New() *Parser {
	return &Parser{
		doc:      &ApiDocSpec{},
		groups:   make(map[string]*ApiGroupSpec),
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

func (p *Parser) Parse(searchDirs []string) error {
	for _, searchDir := range searchDirs {
		fmt.Println("search dir:", searchDir)
		packageDir, err := getPkgName(searchDir)
		if err != nil {
			return err
		}
		if err = p.getAllGoFileInfo(packageDir, searchDir); err != nil {
			return err
		}
	}
	if err := p.packages.ParseTypes(); err != nil {
		return err
	}
	if err := rangeFiles(p.packages.files, p.parseApiInfos); err != nil {
		return err
	}
	return nil
}

func (p *Parser) GetApiDoc() *ApiDocSpec {
	return p.doc
}

func (p *Parser) parseApiInfos(fileName string, astFile *ast.File) error {
	//parse group
	var fileGroup string
	comments := strings.Split(astFile.Doc.Text(), "\n")
	if len(comments) > 0 {
		for _, comment := range comments {
			commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
			if len(commentLine) == 0 {
				continue
			}
			attribute := strings.Fields(commentLine)[0]
			lineRemainder, lowerAttribute := strings.TrimSpace(commentLine[len(attribute):]), strings.ToLower(attribute)
			if lowerAttribute == groupAttr {
				fileGroup = lineRemainder
				break
			}
		}
	}
	for _, comment := range astFile.Comments {
		comments := strings.Split(comment.Text(), "\n")
		if isApiGroupComment(comments) {
			if err := p.parseApiGroupInfo(comments); err != nil {
				return err
			}
			continue
		}
	}
	for _, astDescription := range astFile.Decls {
		switch astDeclaration := astDescription.(type) {
		case *ast.FuncDecl:
			if astDeclaration.Doc != nil && astDeclaration.Doc.List != nil {
				comments := strings.Split(astDeclaration.Doc.Text(), "\n")
				if astDeclaration.Name.Name == "main" { //parse service
					if isApiDocComment(comments) {
						if err := p.parseApiDocInfo(comments); err != nil {
							return err
						}
						continue
					}
				}
				if isApiGroupComment(comments) { //parse group, if in func decl
					if err := p.parseApiGroupInfo(comments); err != nil {
						return err
					}
					continue
				}

				if !isApiComment(comments) {
					continue
				}

				//parse apis
				operation := NewOperation(p)
				for _, comment := range comments {
					err := operation.ParseComment(comment, astFile)
					if err != nil {
						return fmt.Errorf("ParseComment error in file %s :%+v", fileName, err)
					}
				}
				operation.ApiSpec.doc = p.doc                         //ptr, for build full url
				if operation.ApiSpec.Group == "" && fileGroup != "" { //use file group
					operation.ApiSpec.Group = fileGroup
				}
				if operation.ApiSpec.Group == "" {
					p.doc.UngroupedApis = append(p.doc.UngroupedApis, &operation.ApiSpec)
				} else {
					if g, ok := p.groups[operation.ApiSpec.Group]; ok {
						g.Apis = append(g.Apis, &operation.ApiSpec)
					} else {
						group := ApiGroupSpec{
							Group:       operation.ApiSpec.Group,
							Title:       operation.ApiSpec.Group,
							Description: "",
						}
						group.Apis = append(group.Apis, &operation.ApiSpec)
						p.groups[operation.ApiSpec.Group] = &group
						p.doc.Groups = append(p.doc.Groups, &group)
					}
				}
				p.doc.TotalCount += 1
			}
		}
	}

	return nil
}

func (p *Parser) parseApiGroupInfo(comments []string) error {
	previousAttribute := ""
	var group ApiGroupSpec
	for line := 0; line < len(comments); line++ {
		commentLine := comments[line]
		attribute := strings.Split(commentLine, " ")[0]
		value := strings.TrimSpace(commentLine[len(attribute):])
		multilineBlock := false
		if previousAttribute == attribute {
			multilineBlock = true
		}
		switch strings.ToLower(attribute) {
		case groupAttr:
			group.Group = strings.ToLower(value)
		case titleAttr:
			group.Title = value
		case descriptionAttr:
			if multilineBlock {
				group.Description += "\n" + value
				continue
			}
			group.Description = value
		case orderAttr:
			if i, err := strconv.Atoi(value); err == nil {
				group.Order = i
			}
		}
	}
	if group.Group == "" {
		return errors.New("error: group ")
	}
	if g, ok := p.groups[group.Group]; ok {
		g.Group = group.Group
		g.Title = group.Title
		g.Description = group.Description
		g.Order = group.Order
	} else {
		p.groups[group.Group] = &group
		p.doc.Groups = append(p.doc.Groups, &group)
	}
	return nil
}

func (p *Parser) parseApiDocInfo(comments []string) error {
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
		case serviceAttr:
			p.doc.Service = value
		case versionAttr:
			p.doc.Version = value
		case titleAttr:
			p.doc.Title = value
		case descriptionAttr:
			if multilineBlock {
				p.doc.Description += "\n" + value
				continue
			}
			p.doc.Description = value
		case baseURLAttr:
			p.doc.BaseURL = value
		}
	}
	return nil
}

func isApiDocComment(comments []string) bool {
	for _, commentLine := range comments {
		attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
		switch attribute {
		case serviceAttr:
			return true
		case apiAttr, successAttr, failureAttr, responseAttr:
			return false
		}
	}
	return false
}

func isApiGroupComment(comments []string) bool {
	isGroup := false
	for _, commentLine := range comments {
		attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
		switch attribute {
		case serviceAttr, apiAttr, successAttr, failureAttr, responseAttr:
			return false
		case groupAttr:
			isGroup = true
		}
	}
	return isGroup
}

func isApiComment(comments []string) bool {
	isApi := false
	for _, commentLine := range comments {
		attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
		switch attribute {
		case apiAttr, successAttr, failureAttr, requestAttr, responseAttr:
			return true
		}
	}
	return isApi
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
func (p *Parser) getAllGoFileInfo(packageDir, searchDir string) error {
	return filepath.Walk(searchDir, func(path string, f os.FileInfo, _ error) error {
		err := p.Skip(path, f)
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

		return p.parseFile(filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(packageDir, relPath)))), path, nil)
	})
}

func (p *Parser) parseFile(packageDir, path string, src interface{}) error {
	if strings.HasSuffix(strings.ToLower(path), "_test.go") || filepath.Ext(path) != ".go" {
		return nil
	}

	// positions are relative to FileSet
	astFile, err := goparser.ParseFile(token.NewFileSet(), path, src, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("ParseFile error:%+v", err)
	}

	err = p.packages.CollectAstFile(packageDir, path, astFile)
	if err != nil {
		return err
	}

	return nil
}

// Skip returns filepath.SkipDir error if match vendor and hidden folder.
func (p *Parser) Skip(path string, f os.FileInfo) error {
	return walkWith(p.excludes)(path, f)
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

func (p *Parser) getTypeSchema(typeName string, file *ast.File, parentSchema *TypeSchema) (*TypeSchema, error) {
	if IsGolangPrimitiveType(typeName) {
		return &TypeSchema{ //root type
			Name:     typeName,
			FullName: typeName,
			Type:     typeName,
			Parent:   parentSchema,
		}, nil
	}

	typeSpecDef := p.packages.FindTypeSpec(typeName, file, true)
	if typeSpecDef == nil {
		return nil, fmt.Errorf("cannot find type definition: %s", typeName)
	}

	schema, err := p.ParseDefinition(typeSpecDef, parentSchema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// ParseDefinition parses given type spec that corresponds to the type under
// given name and package
func (p *Parser) ParseDefinition(typeSpecDef *TypeSpecDef, parentSchema *TypeSchema) (*TypeSchema, error) {
	typeName := typeSpecDef.FullName()
	if parentSchema != nil && parentSchema.isInTypeChain(typeSpecDef) {
		fmt.Printf("Skipping '%s', recursion detected.\n", typeName)
		return &TypeSchema{
			Name:     typeSpecDef.Name(),
			FullName: typeName,
			Type:     OBJECT,
			PkgPath:  typeSpecDef.PkgPath,
			Parent:   parentSchema,
			Comment:  fmt.Sprintf("%s(Recursion...)", strings.TrimSuffix(typeSpecDef.TypeSpec.Comment.Text(), "\n")),
		}, nil
	}

	// fmt.Printf("Generating %s\n", typeName)

	switch expr := typeSpecDef.TypeSpec.Type.(type) {
	// type Foo struct {...}
	case *ast.StructType:
		return p.parseStruct(typeSpecDef, typeSpecDef.File, expr.Fields, parentSchema)
	case *ast.Ident:
		return p.getTypeSchema(expr.Name, typeSpecDef.File, parentSchema)
	case *ast.SelectorExpr:
		if xIdent, ok := expr.X.(*ast.Ident); ok {
			return p.getTypeSchema(fullTypeName(xIdent.Name, expr.Sel.Name), typeSpecDef.File, parentSchema)
		}
	case *ast.MapType:
		if keyIdent, ok := expr.Key.(*ast.Ident); ok {
			if IsGolangPrimitiveType(keyIdent.Name) {
				example := strings.Trim(getFieldExample(keyIdent.Name, nil), "\"") //map key example
				mapSchema := &TypeSchema{
					Type:       OBJECT,
					Properties: map[string]*TypeSchema{},
					Parent:     parentSchema,
				}
				schema, err := p.parseTypeExpr(typeSpecDef.File, expr.Value, mapSchema)
				if err != nil {
					return nil, err
				}
				mapSchema.TagValue = schema.TagValue
				mapSchema.Name = schema.Name
				mapSchema.FullName = fmt.Sprintf("map[%s]%s", keyIdent.Name, schema.FullName)

				schema.Name = example
				schema.TagValue = ""
				mapSchema.Properties[example] = schema
				return mapSchema, nil
			}
		}

	default:
		fmt.Printf("Type definition of type '%T' is not supported yet. Using 'object' instead.\n", typeSpecDef.TypeSpec.Type)
	}

	sch := TypeSchema{
		Name:    typeName,
		Type:    OBJECT,
		PkgPath: typeSpecDef.PkgPath,
	}

	return &sch, nil
}

func (p *Parser) parseTypeExpr(file *ast.File, typeExpr ast.Expr, parentSchema *TypeSchema) (*TypeSchema, error) {
	switch expr := typeExpr.(type) {
	// type Foo interface{}
	case *ast.InterfaceType:
		return &TypeSchema{
			Name:     "",
			Type:     ANY,
			FullName: ANY,
			Parent:   parentSchema,
		}, nil

	// type Foo Baz
	case *ast.Ident:
		return p.getTypeSchema(expr.Name, file, parentSchema)
	// type Foo *Baz
	case *ast.StarExpr:
		return p.parseTypeExpr(file, expr.X, parentSchema)

	// type Foo pkg.Bar
	case *ast.SelectorExpr:
		if xIdent, ok := expr.X.(*ast.Ident); ok {
			return p.getTypeSchema(fullTypeName(xIdent.Name, expr.Sel.Name), file, parentSchema)
		}
	// type Foo []Baz
	case *ast.ArrayType:
		itemSchema, err := p.parseTypeExpr(file, expr.Elt, parentSchema)
		if err != nil {
			return nil, err
		}
		return &TypeSchema{Type: "array", ArraySchema: itemSchema, Parent: parentSchema, Name: itemSchema.Name, FullName: itemSchema.FullName, TagValue: itemSchema.TagValue}, nil
	// type Foo map[string]Bar
	case *ast.MapType:
		if keyIdent, ok := expr.Key.(*ast.Ident); ok {
			if IsGolangPrimitiveType(keyIdent.Name) {
				example := strings.Trim(getFieldExample(keyIdent.Name, nil), "\"") //map key example
				mapSchema := &TypeSchema{
					Type:       OBJECT,
					Properties: map[string]*TypeSchema{},
					Parent:     parentSchema,
				}
				schema, err := p.parseTypeExpr(file, expr.Value, mapSchema)
				if err != nil {
					return nil, err
				}
				mapSchema.TagValue = schema.TagValue
				mapSchema.Name = schema.Name
				fullName := schema.FullName
				if schema.Type == ARRAY {
					fullName = fmt.Sprintf("[]%s", schema.FullName)
				}
				mapSchema.FullName = fmt.Sprintf("map[%s]%s", keyIdent.Name, fullName)

				schema.Name = example
				schema.TagValue = ""
				mapSchema.Properties[example] = schema
				return mapSchema, nil
			} else {
				return nil, fmt.Errorf("error: map key type %s, just support string or int", keyIdent.Name)
			}
		}
	case *ast.FuncType:
		return nil, errors.New("filed type can't be func")
	// ...
	default:
		fmt.Printf("Type definition of type '%T' is not supported yet. Using 'object' instead.\n", typeExpr)
	}

	return &TypeSchema{Type: OBJECT}, nil
}

func (p *Parser) parseStruct(typeSpecDef *TypeSpecDef, file *ast.File, fields *ast.FieldList, parentSchama *TypeSchema) (*TypeSchema, error) {
	structSchema := &TypeSchema{
		Name:        typeSpecDef.Name(),
		FullName:    typeSpecDef.FullName(),
		PkgPath:     typeSpecDef.PkgPath,
		FullPath:    typeSpecDef.FullPath(),
		Type:        OBJECT,
		Comment:     strings.TrimSuffix(typeSpecDef.TypeSpec.Comment.Text(), "\n"),
		typeSpecDef: typeSpecDef,
		Parent:      parentSchama,
		Properties:  map[string]*TypeSchema{},
	}

	for _, field := range fields.List {
		if field.Names != nil {
			name := field.Names[0].Name
			if !ast.IsExported(name) {
				continue
			}
		}
		schema, err := p.parseTypeExpr(file, field.Type, structSchema)
		if err != nil {
			return nil, err
		}
		if schema == nil {
			continue
		}
		schema.TagValue = getAllTagValue(field)
		schema.Comment = strings.TrimSuffix(field.Comment.Text(), "\n")
		if field.Names == nil { //nested struct, replace with child properties
			for _, p := range schema.Properties {
				if _, ok := structSchema.Properties[strings.ToLower(p.Name)]; !ok { //if not exists key
					structSchema.Properties[strings.ToLower(p.Name)] = p
				}
			}
		} else {
			schema.Name = field.Names[0].Name
			structSchema.Properties[strings.ToLower(schema.Name)] = schema
		}
	}
	return structSchema, nil
}
