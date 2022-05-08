package apidoc

import (
	"fmt"
	"go/ast"
	"reflect"
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
	// Parameters
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
	PkgPath     string
	Required    bool
	Comment     string
	Example     string //example value
	ArraySchema *TypeSchema
	Properties  map[string]*TypeSchema //object
}

func (s *TypeSchema) JSON() string {
	depth := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("// %s %s %s\n", s.Type, s.Name, s.Comment))
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
		if s.Comment != "" {
			sb.WriteString("  //" + buildComment(*s))
		}
		sb.WriteString("\n")
		var i int = 0
		prefix2 := prefix + "  "
		for k, v := range s.Properties {
			sb.WriteString(fmt.Sprintf(prefix2+"\"%s\": ", k)) //write key
			v.parseJSON(depth+1, sb, false)
			haxNext := i < len(s.Properties)-1
			if haxNext {
				sb.WriteString(",")
			}
			//comment
			if len(v.Properties) == 0 && v.Comment != "" && v.ArraySchema == nil {
				sb.WriteString(fmt.Sprintf("  // %s", buildComment(*v)))
			}
			sb.WriteString("\n")
			i++
		}
		sb.WriteString(prefix + "}")

	} else if s.Type == "array" && s.ArraySchema != nil {
		if isNewLine {
			sb.WriteString(prefix + "[")
		} else {
			sb.WriteString("[")
		}
		if s.Comment != "" {
			sb.WriteString(fmt.Sprintf("  // %s", buildComment(*s)))
		}
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
	s := v.Type
	if v.Required {
		s += ", required"
	}
	if len(v.Comment) > 0 {
		s += ", " + v.Comment
	}
	return strings.TrimSuffix(s, "\n")
}

func getExampleValue(typeName string, field *ast.Field) string {
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

	switch typeName {
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint32", "uint64", "byte":
		return fmt.Sprintf("%d", exampleInt(example))
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
		return 0
	}
	v, _ := strconv.Atoi(example)
	return v
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
		return "example"
	}
	return example
}

//getFieldName format json/xml
func getFieldName(name string, field *ast.Field, format string) string {
	if field != nil {
		if field.Tag != nil && field.Tag.Value != "" {
			tag := reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
			if val, ok := tag.Lookup(format); ok {
				name = strings.Split(val, ",")[0]
			}
		}
	}
	return name
}
