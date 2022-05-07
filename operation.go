package apidoc

import (
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"
)

var (
	// 200 Response{data=Data} examples
	responsePattern = regexp.MustCompile(`^(\d+)\s+([\w\-.\\{}=,\[\]]+)\s+(.*)?`)
	// responsePattern = regexp.MustCompile(`^([\w,]+)\s+([\w{}]+)\s+([\w\-.\\{}=,\[\]]+)[^"]*(.*)?`)

	// ResponseType{data1=Type1,data2=Type2}.
	combinedPattern = regexp.MustCompile(`^([\w\-./\[\]]+){(.*)}$`)
)

// Operation describes a single API operation on a path.
type Operation struct {
	parser *Parser
	ApiSpec
}

// NewOperation creates a new Operation with default properties.
// map[int]Response.
func NewOperation(parser *Parser, options ...func(*Operation)) *Operation {
	if parser == nil {
		parser = New()
	}

	result := &Operation{
		parser: parser,
	}

	for _, option := range options {
		option(result)
	}

	return result
}

func (operation *Operation) ParseComment(comment string, astFile *ast.File) error {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	attribute := strings.Fields(commentLine)[0]
	lineRemainder, lowerAttribute := strings.TrimSpace(commentLine[len(attribute):]), strings.ToLower(attribute)

	switch lowerAttribute {
	case titleAttr:
		operation.Title = lineRemainder
	case authorAttr:
		operation.Author = lineRemainder
	case groupAttr:
		operation.Group = lineRemainder
	case acceptAttr:
		operation.Accept = lineRemainder
	case descriptionAttr:
		operation.ParseDescriptionComment(lineRemainder)
	case tagsAttr:
		operation.ParseTagsComment(lineRemainder)
	case apiAttr:
		return operation.ParseRouterComment(lineRemainder)
	case successAttr, failureAttr, responseAttr:
		return operation.ParseResponseComment(lineRemainder, astFile)
	case deprecatedAttr:
		operation.Deprecated = true
	}

	return nil
}

// ParseDescriptionComment godoc.
func (operation *Operation) ParseDescriptionComment(lineRemainder string) {
	if operation.Description == "" {
		operation.Description = lineRemainder

		return
	}
	operation.Description += "\n" + lineRemainder
}

// ParseTagsComment parses comment for given `tag` comment string.
func (operation *Operation) ParseTagsComment(commentLine string) {
	for _, tag := range strings.Split(commentLine, ",") {
		operation.Tags = append(operation.Tags, strings.TrimSpace(tag))
	}
}

var routerPattern = regexp.MustCompile(`^(\w+)[[:blank:]](/[\w./\-{}+:$]*)`)

// ParseRouterComment parses comment for given `router` comment string.
func (operation *Operation) ParseRouterComment(commentLine string) error {
	matches := routerPattern.FindStringSubmatch(commentLine)
	if len(matches) != 3 {
		return fmt.Errorf("can not parse router comment \"%s\"", commentLine)
	}

	httpMethod := strings.ToUpper(matches[1])

	if _, ok := allMethod[httpMethod]; !ok {
		return fmt.Errorf("invalid method: %s", httpMethod)
	}

	operation.HTTPMethod = httpMethod
	operation.Api = matches[2]
	operation.FullURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(operation.parser.doc.BaseURL, "/"), strings.TrimPrefix(operation.Api, "/"))
	return nil
}

// ParseResponseComment parses comment for given `response` comment string.
func (operation *Operation) ParseResponseComment(commentLine string, astFile *ast.File) error {
	matches := responsePattern.FindStringSubmatch(commentLine)
	for i, m := range matches {
		fmt.Println(i, m)
	}
	if len(matches) != 4 {
		return nil
	}
	//200 Response{data=TestData}
	description := strings.Trim(matches[3], "\"")
	codeStr := matches[1]
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return fmt.Errorf("can not parse response comment \"%s\"", commentLine)
	}
	refType := matches[2] //Response{data=TestData}
	//object
	schema, err := operation.parseObject(refType, astFile)
	if err != nil {
		return err
	}
	j := schema.JSON()
	fmt.Printf("schema:%+v\n", schema)
	fmt.Println("json:")
	fmt.Println(j)
	operation.Responses = append(operation.Responses, &ApiResponseSpec{
		StatusCode:  code,
		Format:      "json",
		Examples:    j,
		Schema:      schema,
		Description: description,
	})
	return nil
}

func (operation *Operation) parseObject(refType string, astFile *ast.File) (*TypeSchema, error) {
	switch {
	case IsGolangPrimitiveType(refType):
		refType = TransToValidSchemeType(refType)
		return &TypeSchema{Name: refType, Type: refType}, nil
	case strings.Contains(refType, "{"):
		return operation.parseCombinedObject(refType, astFile)
	default:
		schema, err := operation.parser.getTypeSchema(refType, astFile, nil, true)
		if err != nil {
			return nil, err
		}

		return schema, nil
	}
}

func (operation *Operation) parseCombinedObject(refType string, astFile *ast.File) (*TypeSchema, error) {
	matches := combinedPattern.FindStringSubmatch(refType)
	if len(matches) != 3 { //[Response{data=TestData} Response data=TestData]
		return nil, fmt.Errorf("invalid type: %s", refType)
	}
	fmt.Println("parseCombinedObject matches:", matches)

	schema, err := operation.parseObject(matches[1], astFile)
	if err != nil {
		return nil, err
	}
	fmt.Printf("parseCombinedObject schema=%+v\n", schema)

	fields, props := parseFields(matches[2]), map[string]TypeSchema{}

	for _, field := range fields {
		keyVal := strings.SplitN(field, "=", 2)
		if len(keyVal) == 2 {
			schema, err := operation.parseObject(keyVal[1], astFile)
			if err != nil {
				return nil, err
			}

			props[keyVal[0]] = *schema
		}
	}
	fmt.Printf("props: %+v\n", props)
	if len(props) == 0 {
		return schema, nil
	}
	return schema, nil

	// return spec.ComposedSchema(*schema, spec.Schema{
	// 	SchemaProps: spec.SchemaProps{
	// 		Type:       []string{OBJECT},
	// 		Properties: props,
	// 	},
	// }), nil
}

func parseFields(s string) []string {
	nestLevel := 0
	return strings.FieldsFunc(s, func(char rune) bool {
		if char == '{' {
			nestLevel++

			return false
		} else if char == '}' {
			nestLevel--

			return false
		}
		return char == ',' && nestLevel == 0
	})
}
