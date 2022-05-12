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
	responsePattern = regexp.MustCompile(`^(\d+)\s+([\w\.\d_]+\{.*\}|[\w\.\d_\[\]]+)[^"]*(.*)?`)
	requestPattern  = regexp.MustCompile(`([\w\-.\\\[\]]+)\s*(.*)?`)
	// ResponseType{data1=Type1,data2=Type2}.
	combinedPattern = regexp.MustCompile(`^([\w\-./\[\]]+){(.*)}$`)
	// requestId string true "comment"
	paramPattern = regexp.MustCompile(`(\S+)\s+([\S.]+)\s+(\w+)\s+"([^"]+)"`)
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
	case apiAttr:
		return operation.ParseRouterComment(lineRemainder)
	case requestAttr:
		return operation.ParseRequestComment(lineRemainder, astFile)
	case successAttr, failureAttr, responseAttr:
		return operation.ParseResponseComment(lineRemainder, astFile)
	case headerAttr, queryAttr, paramAttr, formAttr:
		return operation.ParseParametersComment(strings.TrimPrefix(lowerAttribute, "@"), lineRemainder, astFile)
	case deprecatedAttr, "deprecated:":
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

func (operation *Operation) ParseRequestComment(commentLine string, astFile *ast.File) error {
	matches := requestPattern.FindStringSubmatch(commentLine)
	//0 Request 1 Request 2 Comment
	if len(matches) != 3 {
		return nil
	}
	refType := matches[1]
	switch {
	case IsGolangPrimitiveType(refType):
		return nil
	default:
		schema, err := operation.parser.getTypeSchema(refType, astFile, nil, nil)
		if err != nil {
			return err
		}
		if schema == nil || schema.Properties == nil {
			return nil
		}
		operation.Requests = ApiRequestSpec{
			Parameters: map[string]*ApiParameterSpec{},
		}
		var parameterCount = 0
		for _, p := range schema.Properties {
			tags := p.ParameterTags()
			if tags != nil {
				if len(tags) > 0 && !p.HasJSONTag() {
					parameterCount++
				}
			}
			for pType, pName := range tags {
				if param, ok := operation.Requests.Parameters[pName]; ok {
					param.parameterTypes = append(param.parameterTypes, pType)
				} else {
					operation.Requests.Parameters[pName] = &ApiParameterSpec{
						Name:           pName,
						Required:       p.IsRequired(),
						Description:    p.Comment,
						Validate:       p.ValidateTag(),
						parameterTypes: []string{pType},
						DataType:       p.Type,
					}
				}
			}
		}
		if parameterCount < len(schema.Properties) {
			operation.Requests.Body = schema.JSON()
		}
		return nil
	}
}

//ParseParametersComment parses parameters (@header, @param, @query, @form)
//@param [name] [type] [required] [comment]
//@query demo int true "测试参数"
func (operation *Operation) ParseParametersComment(parameterType, commentLine string, astFile *ast.File) error {
	matches := paramPattern.FindStringSubmatch(commentLine)
	if len(matches) != 5 {
		return fmt.Errorf("missing required param comment parameters \"%s\"", commentLine)
	}
	name := matches[1]
	dataType := matches[2]
	required := strings.ToLower(matches[3]) == "true"
	description := matches[4]
	if _, ok := operation.Requests.Parameters[name]; !ok {
		operation.Requests.Parameters[name] = &ApiParameterSpec{
			Name:           name,
			DataType:       dataType,
			Required:       required,
			Description:    description,
			parameterTypes: []string{parameterType},
		}
	}
	return nil
}

// ParseResponseComment parses comment for given `response` comment string.
func (operation *Operation) ParseResponseComment(commentLine string, astFile *ast.File) error {
	matches := responsePattern.FindStringSubmatch(commentLine)
	if len(matches) != 4 && len(matches) != 3 {
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
	// fmt.Printf("schema:%+v\n", schema)
	// fmt.Println("json:")
	// fmt.Println(j)
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
	arrayFlag := "[]"
	isArray := strings.HasPrefix(refType, arrayFlag)
	if isArray { //array
		typeName := strings.TrimPrefix(refType, arrayFlag)
		schema, err := operation.parseObject(typeName, astFile)
		if err != nil {
			return nil, err
		}
		return &TypeSchema{
			Name:        refType,
			Type:        ARRAY,
			ArraySchema: schema,
		}, nil
	}
	switch {
	case IsGolangPrimitiveType(refType):
		return &TypeSchema{
			Name:     refType,
			Type:     refType,
			FullName: refType,
		}, nil
	case strings.Contains(refType, "{"):
		return operation.parseCombinedObject(refType, astFile)
	default:
		schema, err := operation.parser.getTypeSchema(refType, astFile, nil, nil)
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

	schemaA, err := operation.parseObject(matches[1], astFile)
	if err != nil {
		return nil, err
	}

	fields := parseFields(matches[2])
	// props := map[string]TypeSchema{}
	for _, field := range fields {
		keyVal := strings.SplitN(field, "=", 2)
		if len(keyVal) == 2 {
			// fmt.Println("keyVal", keyVal[0], keyVal[1]) //data TestData
			// if is number or string wrap, replace it
			if isReplaceValue(keyVal[1]) { //replace int,string, examples code or msg
				if p, ok := schemaA.Properties[keyVal[0]]; ok {
					p.example = keyVal[1] //replace response code, msg
				}
			} else {
				//check is array
				arrayFlag := "[]"
				typeName := keyVal[1]
				isArray := strings.HasPrefix(typeName, arrayFlag)
				if isArray { //array
					typeName = strings.TrimPrefix(typeName, arrayFlag)
				}
				schema, err := operation.parseObject(typeName, astFile)
				if err != nil {
					return nil, err
				}
				key := strings.ToLower(keyVal[0])
				if isArray {
					arrSchema := &TypeSchema{
						Name:        key,
						Type:        ARRAY,
						ArraySchema: schema,
					}
					if old, ok := schemaA.Properties[key]; ok {
						arrSchema.TagValue = old.TagValue // use old tag, for example Response.data
					}
					schemaA.Properties[key] = arrSchema
				} else {
					if old, ok := schemaA.Properties[key]; ok {
						schema.TagValue = old.TagValue // use old tag, for example Response.data
					}
					schemaA.Properties[key] = schema //data=xx
				}
			}

		}
	}
	return schemaA, nil
}

func isReplaceValue(val string) bool {
	if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) || (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
		return true
	}
	_, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return true
	}
	_, err = strconv.ParseFloat(val, 64)
	if err == nil {
		return true
	}
	_, err = strconv.ParseBool(val)
	return err == nil
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
