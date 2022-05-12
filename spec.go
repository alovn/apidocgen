package apidoc

import (
	"fmt"
	"go/ast"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type ApiDocSpec struct {
	Service     string
	Title       string
	Version     string
	Description string
	Scheme      string
	BaseURL     string
	Groups      []*ApiGroupSpec
	Apis        []*ApiSpec
}

type ApiGroupSpec struct {
	Group       string
	Title       string
	Description string
	Apis        []*ApiSpec
}

type ApiSpec struct {
	Title       string
	HTTPMethod  string
	Api         string
	FullURL     string
	Version     string
	Accept      string //json,xml,form
	Format      string //json,xml
	Description string
	Author      string
	Deprecated  bool
	Group       string
	Responses   []*ApiResponseSpec
	Requests    ApiRequestSpec
}

type ApiRequestSpec struct {
	Parameters map[string]*ApiParameterSpec
	Body       string
}

type ApiParameterSpec struct {
	Name           string
	DataType       string
	Required       bool
	Description    string
	Validate       string
	Example        string
	parameterTypes []string
}

func (p ApiParameterSpec) ParameterTypes() string {
	s := ""
	for i, m := range p.parameterTypes {
		if i == 0 {
			s += m
		} else {
			s += "," + m
		}
	}
	return s
}

type ApiResponseSpec struct {
	StatusCode  int
	Format      string //json xml
	Examples    string
	Schema      *TypeSchema
	Description string
}

type TypeSchema struct {
	Name          string //xxRequest, xxResponse
	Type          string //int, string, bool, object, array, any
	FullName      string
	PkgPath       string
	Comment       string
	ArraySchema   *TypeSchema
	Properties    map[string]*TypeSchema //object
	Parent        *TypeSchema
	TagValue      string
	typeSpecDef   *TypeSpecDef
	example       string //example value
	parameterTags map[string]string
}

func (s *TypeSchema) JSON() string {
	depth := 0
	var sb strings.Builder
	// sb.WriteString(fmt.Sprintf("// %s %s %s\n", s.Type, s.Name, s.Comment))
	s.parseJSON(depth, &sb, true)
	return sb.String()
}

func (s *TypeSchema) parseJSON(depth int, sb *strings.Builder, isNewLine bool) {
	prefix := ""
	for i := depth; i > 0; i-- {
		prefix += "  "
	}

	if s.Type == OBJECT && s.Properties != nil && len(s.Properties) > 0 {
		if isNewLine {
			sb.WriteString(prefix + "{")
		} else {
			sb.WriteString("{")
		}
		sb.WriteString("  //" + s.buildComment())
		sb.WriteString("\n")
		var i int = 0
		prefix2 := prefix + "  "
		//sort keys
		var keys []string
		for k := range s.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		//build json
		for _, k := range keys {
			v := s.Properties[k]
			key, isOmitempty := v.JSONKey()
			if isOmitempty && v.ExampleValue() == NULL { //omitempty
				continue
			}
			if v.parameterTags != nil && !v.HasJSONTag() {
				continue
			}
			sb.WriteString(fmt.Sprintf(prefix2+"\"%s\": ", key)) //write key
			v.parseJSON(depth+1, sb, false)
			haxNext := i < len(s.Properties)-1
			if haxNext {
				sb.WriteString(",")
			}
			//comment
			if len(v.Properties) == 0 && v.ArraySchema == nil {
				sb.WriteString(fmt.Sprintf("  //%s", v.buildComment()))
			}
			sb.WriteString("\n")
			i++
		}
		sb.WriteString(prefix + "}")

	} else if s.Type == ARRAY && s.ArraySchema != nil {
		if isNewLine {
			sb.WriteString(prefix + "[")
		} else {
			sb.WriteString("[")
		}
		sb.WriteString(fmt.Sprintf("  //%s", s.buildComment()))
		sb.WriteString("\n")
		if s.ArraySchema.ExampleValue() != NULL {
			s.ArraySchema.parseJSON(depth+1, sb, true)
		}
		sb.WriteString("\n")
		sb.WriteString(prefix + "]")
	} else { // write example value
		if depth == 0 {
			sb.WriteString(fmt.Sprintf("//%s\n", s.buildComment()))
		}
		if isNewLine {
			sb.WriteString(prefix + s.ExampleValue())
		} else {
			sb.WriteString(s.ExampleValue())
		}
	}
}

func (v *TypeSchema) buildComment() string {
	if v == nil {
		return ""
	}
	s := ""
	if v.Type == ARRAY {
		arrayName := v.ArraySchema.Type //int
		if v.ArraySchema.Type == OBJECT {
			arrayName = v.ArraySchema.FullName
		}
		s += fmt.Sprintf("%s[%s]", ARRAY, arrayName)
	} else if len(v.Properties) > 0 { //object
		s += fmt.Sprintf("%s(%s)", v.Type, v.PkgPath+v.FullName)
	} else {
		s += v.Type
	}
	if v.IsRequired() {
		s += ", required"
	}
	if len(v.Comment) > 0 {
		s += ", " + v.Comment
	}
	return strings.TrimSuffix(s, "\n")
}

func (v *TypeSchema) isInTypeChain(typeSpecDef *TypeSpecDef) bool {
	if v.typeSpecDef != nil {
		if v.typeSpecDef == typeSpecDef {
			return true
		}
	}
	if v.Parent != nil {
		return v.Parent.isInTypeChain(typeSpecDef)
	}
	return false
}

func (v *TypeSchema) JSONKey() (key string, isOmitempty bool) {
	if v.TagValue == "" {
		return v.Name, false
	}
	val, has := v.GetTag("json")
	if has {
		key = strings.Split(val, ",")[0]
		isOmitempty = strings.Contains(val, "omitempty")
		return
	} else {
		key = v.Name
		isOmitempty = false
		return
	}
}

func (v *TypeSchema) HasJSONTag() bool {
	_, has := v.GetTag("json")
	return has
}

func (v *TypeSchema) ExampleValue() string {
	if v.Type == ARRAY && v.ArraySchema == nil {
		return NULL
	}
	if v.Type == OBJECT && (v.Properties == nil || len(v.Properties) == 0) {
		return NULL
	}
	if v.example != "" {
		return v.example
	}
	example := ""

	if val, has := v.GetTag("example"); has {
		example = val
	}
	v.example = getTypeExample(v.Type, example)
	return v.example
}

func (v *TypeSchema) GetTag(name string) (value string, has bool) {
	tag := reflect.StructTag(v.TagValue)
	return tag.Lookup(name)
}

func (v *TypeSchema) ParameterTags() map[string]string {
	if v.parameterTags != nil {
		return v.parameterTags
	}
	parameterTags := make(map[string]string)
	keys := []string{"header", "param", "query", "form"}
	for _, key := range keys {
		if val, has := v.GetTag(key); has {
			parameterTags[key] = val
		}
	}
	v.parameterTags = parameterTags
	return v.parameterTags
}

func (v *TypeSchema) IsRequired() (required bool) {
	if val, has := v.GetTag("required"); has {
		required, _ = strconv.ParseBool(val)
	}
	if !required {
		validate := v.ValidateTag()
		required = strings.Contains(validate, "required")
	}
	return
}

func (v *TypeSchema) ValidateTag() (validate string) {
	validate, _ = v.GetTag("validate")
	return
}

func getFieldExample(typeName string, field *ast.Field) string {
	example := ""
	if field != nil {
		name := field.Names[0]
		if !name.IsExported() {
			return ""
		}
		if field.Tag != nil && field.Tag.Value != "" {
			tag := reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
			if val, ok := tag.Lookup("example"); ok {
				example = val
			}
		}
	}
	return getTypeExample(typeName, example)
}

func getTypeExample(typeName, example string) string {
	switch typeName {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		return fmt.Sprintf("%d", exampleInt(example))
	case "float32", "float64":
		return fmt.Sprintf("%.2f", exampleFloat(example))
	case "rune":
		return fmt.Sprintf("'%c'", exampleRune(example))
	case "string":
		return fmt.Sprintf("\"%s\"", exampleString(example))
	case "bool":
		return fmt.Sprintf("%t", exampleBool(example))
	case "any":
		return "null"
	}

	return example
}

func exampleInt(example string) int {
	if example == "" {
		return 123
	}
	v, _ := strconv.Atoi(example)
	return v
}

func exampleFloat(example string) float64 {
	val := 1.23
	if example != "" {
		if v, err := strconv.ParseFloat(example, 64); err == nil {
			val = v
		}
	}
	return val
}
func exampleRune(example string) rune {
	val := rune(97)
	if example == "" {
		return val
	}
	for _, r := range example {
		if r == ' ' {
			continue
		}
		return r
	}
	return val
}

func exampleBool(example string) bool {
	val := true
	if example == "" {
		val, _ = strconv.ParseBool(example)
	}
	return val
}

func exampleString(example string) string {
	if example == "" {
		return "abc"
	}
	return example
}

func getAllTagValue(field *ast.Field) string {
	if field != nil && field.Tag != nil && field.Tag.Value != "" {
		return strings.ReplaceAll(field.Tag.Value, "`", "")
	}
	return ""
}
