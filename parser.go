package apidoc

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
	"strings"
)

const (
	apiAttr         = "@api"
	titleAttr       = "@title"
	groupAttr       = "@group"
	versionAttr     = "@version"
	descriptionAttr = "@desc"
	acceptAttr      = "@accept"
	successAttr     = "@success"
	failureAttr     = "@failure"
	responseAttr    = "@response"
	deprecatedAttr  = "@deprecated"
	tagsAttr        = "@tags"
	authorAttr      = "@author"

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
}

type Parser struct {
	doc      *ApiDocSpec
	groups   map[string]*ApiGroupSpec
	packages *PackagesDefinitions
	// excludes excludes dirs and files in SearchDir
	excludes map[string]struct{}
	// structStack stores full names of the structures that were already parsed or are being parsed now
	structStack []*TypeSpecDef
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
	if err = parser.packages.ParseTypes(); err != nil {
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
		if isApiGroupComment(comments) {
			if err := parser.parseApiGroupInfo(comments); err != nil {
				return err
			}
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
			if operation.ApiSpec.Group == "" {
				parser.doc.Apis = append(parser.doc.Apis, &operation.ApiSpec)
			} else {
				if g, ok := parser.groups[operation.ApiSpec.Group]; ok {
					g.Apis = append(g.Apis, &operation.ApiSpec)
				} else {
					group := ApiGroupSpec{
						Group:       operation.ApiSpec.Group,
						Title:       operation.ApiSpec.Group,
						Description: "",
					}
					group.Apis = append(group.Apis, &operation.ApiSpec)
					parser.groups[operation.ApiSpec.Group] = &group
					parser.doc.Groups = append(parser.doc.Groups, &group)
				}
			}
		}
	}

	return nil
}

func (parser *Parser) parseApiGroupInfo(comments []string) error {
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
		}
	}
	if group.Group == "" {
		return errors.New("error: group ")
	}
	if g, ok := parser.groups[group.Group]; ok {
		g.Group = group.Group
		g.Title = group.Title
		g.Description = group.Description
	} else {
		parser.groups[group.Group] = &group
		parser.doc.Groups = append(parser.doc.Groups, &group)
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
		case baseURLAttr:
			parser.doc.BaseURL = value
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

func isApiGroupComment(comments []string) bool {
	for _, commentLine := range comments {
		attribute := strings.ToLower(strings.Split(commentLine, " ")[0])
		switch attribute {
		case apiAttr, successAttr, failureAttr, responseAttr:
			return false
		case groupAttr:
			return true
		}
	}
	return false
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

func (parser *Parser) getTypeSchema(typeName string, file *ast.File, field *ast.Field, ref bool) (*TypeSchema, error) {
	if IsGolangPrimitiveType(typeName) {
		name := field.Names[0].Name
		fieldName := getFieldName(name, field, "json")
		return &TypeSchema{
			Name:      name,
			FieldName: fieldName,
			Comment:   field.Comment.Text(),
			Type:      typeName,
			Example:   getExampleValue(typeName, field),
		}, nil
	}

	typeSpecDef := parser.packages.FindTypeSpec(typeName, file, true)
	if typeSpecDef == nil {
		return nil, fmt.Errorf("cannot find type definition: %s", typeName)
	}
	fmt.Println("typeSpecDef", typeSpecDef.Name())

	schema, err := parser.ParseDefinition(typeSpecDef)
	if err != nil {
		return nil, err
	}

	// if ref && len(schema.Type) > 0 && schema.Type[0] == OBJECT {
	// 	return parser.getRefTypeSchema(typeSpecDef, schema), nil
	// }

	return schema, nil
}

func (parser *Parser) isInStructStack(typeSpecDef *TypeSpecDef) bool {
	for _, specDef := range parser.structStack {
		if typeSpecDef == specDef {
			return true
		}
	}

	return false
}

// ParseDefinition parses given type spec that corresponds to the type under
// given name and package
func (parser *Parser) ParseDefinition(typeSpecDef *TypeSpecDef) (*TypeSchema, error) {
	typeName := typeSpecDef.FullName()
	refTypeName := TypeDocName(typeName, typeSpecDef.TypeSpec)

	if parser.isInStructStack(typeSpecDef) {
		fmt.Printf("Skipping '%s', recursion detected.", typeName)
		return &TypeSchema{
			Name:    refTypeName,
			Type:    OBJECT,
			PkgPath: typeSpecDef.PkgPath,
		}, nil
	}

	parser.structStack = append(parser.structStack, typeSpecDef)

	fmt.Printf("Generating %s\n", typeName)

	switch expr := typeSpecDef.TypeSpec.Type.(type) {
	// type Foo struct {...}
	case *ast.StructType:
		return parser.parseStruct(typeSpecDef.File, expr.Fields)
	default:
		fmt.Printf("Type definition of type '%T' is not supported yet. Using 'object' instead.\n", typeSpecDef.TypeSpec.Type)
	}

	sch := TypeSchema{
		Name:    refTypeName,
		Type:    OBJECT,
		PkgPath: typeSpecDef.PkgPath,
	}

	return &sch, nil
}

func (parser *Parser) parseTypeExpr(file *ast.File, field *ast.Field, typeExpr ast.Expr, ref bool) (*TypeSchema, error) {
	switch expr := typeExpr.(type) {
	// type Foo interface{}
	case *ast.InterfaceType:
		return &TypeSchema{}, nil

	// type Foo struct {...}
	case *ast.StructType:
		return parser.parseStruct(file, expr.Fields)

	// type Foo Baz
	case *ast.Ident:
		return parser.getTypeSchema(expr.Name, file, field, ref)

	// type Foo *Baz
	case *ast.StarExpr:
		return parser.parseTypeExpr(file, field, expr.X, ref)

	// type Foo pkg.Bar
	case *ast.SelectorExpr:
		if xIdent, ok := expr.X.(*ast.Ident); ok {
			return parser.getTypeSchema(fullTypeName(xIdent.Name, expr.Sel.Name), file, field, ref)
		}
	// type Foo []Baz
	case *ast.ArrayType:
		itemSchema, err := parser.parseTypeExpr(file, field, expr.Elt, true)
		if err != nil {
			return nil, err
		}
		return &TypeSchema{Type: "array", ArraySchema: itemSchema}, nil
	// type Foo map[string]Bar
	// case *ast.MapType:
	// 	if _, ok := expr.Value.(*ast.InterfaceType); ok {
	// 		return &TypeSchema{Type: OBJECT, Properties: nil}, nil
	// 	}
	// 	schema, err := parser.parseTypeExpr(file, expr.Value, true)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return spec.MapProperty(schema), nil

	// case *ast.FuncType:
	// 	return nil, ErrFuncTypeField
	// ...
	default:
		fmt.Printf("Type definition of type '%T' is not supported yet. Using 'object' instead.\n", typeExpr)
	}

	return &TypeSchema{Type: OBJECT}, nil
}

func (parser *Parser) parseStruct(file *ast.File, fields *ast.FieldList) (*TypeSchema, error) {
	properties := make(map[string]*TypeSchema)

	for _, field := range fields.List {
		if len(field.Names) != 1 {
			return nil, errors.New("error len(field.Names) != 1")
		}
		// name := field.Names[0].Name
		schema, err := parser.parseStructField(file, field)
		if err != nil {
			return nil, err
		}
		properties[schema.FieldName] = schema

		// name := field.Names[0]
		// key := name.Name
		// fmt.Printf("%s %v\n\n\n", name, name.Obj.Decl)
		// if field.Tag != nil && field.Tag.Value != "" {
		// 	tag := reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
		// 	if j, ok := tag.Lookup("json"); ok && j != "" { //xx,omitempty
		// 		key = strings.Split(j, ",")[0]
		// 	}
		// }
		// fmt.Println(key)

		// example, ok := tag.Lookup("example")
		// if ok {
		// 	m[fieldKey] = example
		// }
		// fmt.Println("field", field)
	}
	return &TypeSchema{
		Name:       file.Name.Name,
		Type:       OBJECT,
		Properties: properties,
	}, nil
}

func (parser *Parser) parseStructField(file *ast.File, field *ast.Field) (*TypeSchema, error) {
	name := field.Names[0].Name
	if !ast.IsExported(name) {
		return nil, nil
	}

	typeName, err := getFieldType(field.Type)
	if err != nil {
		return nil, err
	}

	schema, err := parser.getTypeSchema(typeName, file, field, false)
	if err != nil {
		return nil, err
	}

	return schema, nil

	// if field.Names == nil {
	// 	typeName, err := getFieldType(field.Type)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	schema, err := parser.getTypeSchema(typeName, file, false)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if len(schema.Type) > 0 && schema.Type == OBJECT {
	// 		if len(schema.Properties) == 0 {
	// 			return nil, nil
	// 		}

	// 		properties := map[string]*TypeSchema{}
	// 		for k, v := range schema.Properties {
	// 			properties[k] = v
	// 		}

	// 		return properties, nil
	// 	}

	// 	// for alias type of non-struct types ,such as array,map, etc. ignore field tag.
	// 	// return map[string]*TypeSchema{Type: typeName}, nil, nil
	// }

	// ps := parser.fieldParserFactory(parser, field)

	// if ps.ShouldSkip() {
	// 	return nil, nil, nil
	// }

	// fieldName, err := ps.FieldName()
	// if err != nil {
	// 	return nil, nil, err
	// }

	// schema, err := ps.CustomSchema()
	// if err != nil {
	// 	return nil, nil, err
	// }

	// if schema == nil {
	// 	typeName, err := getFieldType(field.Type)
	// 	if err == nil {
	// 		// named type
	// 		schema, err = parser.getTypeSchema(typeName, file, true)
	// 	} else {
	// 		// unnamed type
	// 		schema, err = parser.parseTypeExpr(file, field.Type, false)
	// 	}

	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// }

	// err = ps.ComplementSchema(schema)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// var tagRequired []string

	// required, err := ps.IsRequired()
	// if err != nil {
	// 	return nil, nil, err
	// }

	// if required {
	// 	tagRequired = append(tagRequired, fieldName)
	// }

	// return map[string]*TypeSchema{fieldName: *schema}, tagRequired, nil
}

func getFieldType(field ast.Expr) (string, error) {
	switch fieldType := field.(type) {
	case *ast.Ident:
		return fieldType.Name, nil
	case *ast.SelectorExpr:
		packageName, err := getFieldType(fieldType.X)
		if err != nil {
			return "", err
		}

		return fullTypeName(packageName, fieldType.Sel.Name), nil
	case *ast.StarExpr:
		fullName, err := getFieldType(fieldType.X)
		if err != nil {
			return "", err
		}

		return fullName, nil
	case *ast.InterfaceType:
		return ANY, nil
	default:
		return "", fmt.Errorf("unknown field type %#v", field)
	}
}
