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
	Tags        []string
	Group       string
	Responses   []*ApiResponseSpec
	Requests    ApiRequestSpec
}

type ApiRequestSpec struct {
	Parameters map[string]*ApiParameterSpec
	Body       string
}

type ApiParameterSpec struct {
	Name        string
	Required    bool
	Description string
	Validate    string
	Example     string
	types       []string
}

func (p ApiParameterSpec) Types() string {
	s := ""
	for i, t := range p.types {
		if i == 0 {
			s += t
		} else {
			s += "," + t
		}
	}
	return s
}

type ApiResponseSpec struct {
	IsSuccess   bool
	StatusCode  int
	Format      string //json xml
	Examples    string
	Schema      *TypeSchema
	Description string
}

type TypeSchema struct {
	Name        string //xxRequest, xxResponse
	FieldName   string
	Type        string //int, string, bool, object, array
	FullName    string
	PkgPath     string
	Required    bool
	Comment     string
	Example     string //example value
	IsArray     bool
	ArraySchema *TypeSchema
	Properties  map[string]*TypeSchema //object
	IsOmitempty bool
	Validate    string
	Tags        map[string]string
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

	if s.Type == OBJECT && s.Properties != nil {
		if isNewLine {
			sb.WriteString(prefix + "{")
		} else {
			sb.WriteString("{")
		}
		sb.WriteString("  //" + buildComment(*s))
		sb.WriteString("\n")
		var i int = 0
		prefix2 := prefix + "  "
		//sort keys
		var keys []string
		for k := range s.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := s.Properties[k]
			if v.IsOmitempty && v.Example == "null" { //omitempty
				continue
			}
			sb.WriteString(fmt.Sprintf(prefix2+"\"%s\": ", k)) //write key
			v.parseJSON(depth+1, sb, false)
			haxNext := i < len(s.Properties)-1
			if haxNext {
				sb.WriteString(",")
			}
			//comment
			if len(v.Properties) == 0 && v.ArraySchema == nil {
				sb.WriteString(fmt.Sprintf("  //%s", buildComment(*v)))
			}
			sb.WriteString("\n")
			i++
		}
		sb.WriteString(prefix + "}")

	} else if s.IsArray && s.Type == ARRAY && s.ArraySchema != nil {
		if isNewLine {
			sb.WriteString(prefix + "[")
		} else {
			sb.WriteString("[")
		}
		sb.WriteString(fmt.Sprintf("  //%s", buildComment(*s)))
		sb.WriteString("\n")
		s.ArraySchema.parseJSON(depth+1, sb, true)
		sb.WriteString("\n")

		sb.WriteString(prefix + "]")
	} else { // write example value
		if isNewLine {
			sb.WriteString(prefix + s.Example)
		} else {
			sb.WriteString(s.Example)
		}
	}
}

func buildComment(v TypeSchema) string {
	s := ""
	if v.IsArray {
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
	if v.Required {
		s += ", required"
	}
	if len(v.Comment) > 0 {
		s += ", " + v.Comment
	}
	return strings.TrimSuffix(s, "\n")
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
	return getExampleValue(typeName, example)
}

func getExampleValue(typeName, example string) string {
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

//getFieldName format json/xml
//getFiledName("abc", field, "json")
func getFieldName(name string, field *ast.Field, format string) (isOmitempty bool, fieldName string) {
	tagValue, has := getTagValue(format, field)
	if has {
		fieldName = strings.Split(tagValue, ",")[0]
		isOmitempty = strings.Contains(tagValue, "omitempty")
		return
	}
	return false, name
}

func getTagValue(tagName string, field *ast.Field) (tagValue string, has bool) {
	if field != nil {
		if field.Tag != nil && field.Tag.Value != "" {
			tag := reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
			return tag.Lookup(tagName)
		}
	}
	return
}

func getValidateTagValue(field *ast.Field) (validate string) {
	validate, _ = getTagValue("validate", field)
	return
}

func getRequiredTagValue(field *ast.Field) (required bool) {
	if val, has := getTagValue("required", field); has {
		required, _ = strconv.ParseBool(val)
	}
	if !required {
		validate := getValidateTagValue(field)
		required = strings.Contains(validate, "required")
	}
	return
}

func getParameterTags(field *ast.Field) map[string]string {
	parameters := make(map[string]string)

	keys := []string{"header", "param", "query", "form"}
	for _, key := range keys {
		if val, has := getTagValue(key, field); has {
			parameters[key] = val
		}
	}
	return parameters
}
