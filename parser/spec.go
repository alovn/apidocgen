package parser

import (
	"fmt"
	"go/ast"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ApiDocSpec struct {
	Service       string
	Title         string
	Version       string
	Description   string
	BaseURL       string
	Groups        []*ApiGroupSpec
	UngroupedApis []*ApiSpec
	TotalCount    int
}

type ApiGroupSpec struct {
	Group       string
	Title       string
	Description string
	Apis        []*ApiSpec
	Order       int // sort
}

type ApiSpec struct {
	doc         *ApiDocSpec
	Title       string
	HTTPMethod  string
	Api         string
	Version     string
	Accept      string // json,xml,form
	Format      string // json,xml
	Description string
	Author      string
	Deprecated  bool
	Group       string
	Responses   []*ApiResponseSpec
	Requests    ApiRequestSpec
	Order       int // sort
}

func (a *ApiSpec) FullURL() string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(a.doc.BaseURL, "/"), strings.TrimPrefix(a.Api, "/"))
}

type ApiRequestSpec struct {
	Parameters map[string]*ApiParameterSpec
	Accept     string // accept format
	Schema     *TypeSchema
	body       string // ache
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
	Format      string // json xml
	Schema      *TypeSchema
	Description string
	IsMock      bool
	body        string // cache
	pureBody    string // cache
}

type TypeSchema struct {
	Name          string // xxRequest, xxResponse, for example: RegisterRequest
	Type          string // int, string, bool, object, array, any
	FullName      string
	PkgPath       string // for example: github.com/alovn/apidocgen/examples/svc-user/handler
	FullPath      string // for example: github.com/alovn/apidocgen/examples/svc-user/handler.RegisterRequest
	Comment       string
	ArraySchema   *TypeSchema
	Properties    map[string]*TypeSchema // object
	Parent        *TypeSchema
	TagValue      string
	typeSpecDef   *TypeSpecDef
	example       string // example value
	parameterTags map[string]string
	xmlName       string
}

func (s *ApiRequestSpec) Body() string {
	if s.Schema == nil {
		return ""
	}
	if s.body != "" {
		return s.body
	}
	s.body = s.Schema.Write(s.Accept, true)
	return s.body
}

func (s *ApiResponseSpec) Body() string {
	if s.Schema == nil {
		return ""
	}
	if s.body != "" {
		return s.body
	}
	s.body = s.Schema.Write(s.Format, true)
	return s.body
}

func (s *ApiResponseSpec) PureBody() string {
	if s.Schema == nil {
		return ""
	}
	if s.pureBody != "" {
		return s.pureBody
	}
	s.pureBody = s.Schema.Write(s.Format, false)
	return s.pureBody
}

func (s *TypeSchema) Write(format string, withComment bool) (body string) {
	switch format {
	case "json":
		return s.JSON(withComment)
	case "jsonp":
		return s.JSONP(withComment)
	case "xml":
		return s.XML(withComment)
	default:
		return s.JSON(withComment)
	}
}

func (s *TypeSchema) JSONP(withComment bool) string {
	var sb strings.Builder

	sb.WriteString("callback(")
	sb.WriteString(s.JSON(withComment))
	sb.WriteString(")")
	return sb.String()
}

func (s *TypeSchema) JSON(withComment bool) string {
	depth := 0
	var sb strings.Builder
	// sb.WriteString(fmt.Sprintf("// %s %s %s\n", s.Type, s.Name, s.Comment))
	s.parseJSON(depth, &sb, false, withComment)
	return sb.String()
}

func (s *TypeSchema) parseJSON(depth int, sw io.StringWriter, isNewLine, withComment bool) {
	prefix := ""
	for i := depth; i > 0; i-- {
		prefix += "  "
	}

	if s.Type == OBJECT && s.Properties != nil && len(s.Properties) > 0 {
		if isNewLine {
			sw.WriteString(prefix)
		}

		sw.WriteString("{" + s.buildComment(withComment) + "\n")
		prefix2 := prefix + "  "
		// sort keys
		var keys []string
		for k, v := range s.Properties {
			if v.isIgnoreJsonKey() {
				continue
			}

			keys = append(keys, k)
		}

		sort.Strings(keys)
		var i int = 0
		// build json
		for _, k := range keys {
			v := s.Properties[k]
			key, _ := v.JSONKey()
			sw.WriteString(fmt.Sprintf(prefix2+"\"%s\": ", key)) // write key
			v.parseJSON(depth+1, sw, false, withComment)
			haxNext := i < len(keys)-1
			if haxNext {
				sw.WriteString(",")
			}
			// comment
			if len(v.Properties) == 0 && v.ArraySchema == nil {
				sw.WriteString(v.buildComment(withComment))
			}

			sw.WriteString("\n")
			i++
		}

		sw.WriteString(prefix + "}")
	} else if s.Type == ARRAY && s.ArraySchema != nil {
		if isNewLine {
			sw.WriteString(prefix)
		}
		sw.WriteString(fmt.Sprintf("[%s\n", s.buildComment(withComment)))
		if s.ArraySchema.ExampleValue() != NULL {
			s.ArraySchema.parseJSON(depth+1, sw, true, withComment)
		}
		sw.WriteString("\n")
		sw.WriteString(prefix + "]")
	} else { // write example value
		if withComment {
			if depth == 0 {
				sw.WriteString(fmt.Sprintf("%s\n", strings.TrimLeft(s.buildComment(withComment), " ")))
			}
		}
		if isNewLine {
			sw.WriteString(prefix)
		}
		sw.WriteString(s.ExampleValue())
	}
}

func (s *TypeSchema) isIgnoreJsonKey() (isIgnore bool) {
	if s.Name == "XMLName" && s.FullName == "xml.Name" { // ignore xml
		isIgnore = true
		return
	}
	key, isOmitempty := s.JSONKey()
	if key == "-" || (isOmitempty && s.ExampleValue() == NULL) { // omitempty
		isIgnore = true
		return
	}
	if s.parameterTags != nil && !s.hasJSONTag() { // request parameter ignore
		isIgnore = true
		return
	}
	return false
}

func (s *TypeSchema) XML(withComment bool) string {
	depth := 0

	var sb strings.Builder
	// sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	// sb.WriteString(fmt.Sprintf("// %s %s %s\n", s.Type, s.Name, s.Comment))
	s.parseXML(depth, &sb, false, withComment)
	return sb.String()
}

func (s *TypeSchema) parseXML(depth int, sw io.StringWriter, isNewLine, withComment bool) {
	prefix := ""
	for i := depth; i > 0; i-- {
		prefix += "  "
	}

	if s.Type == OBJECT && s.Properties != nil && len(s.Properties) > 0 {
		if isNewLine {
			sw.WriteString("\n")
		}
		xmlName := s.XMLName()
		attrs := s.XMLAttrs()

		sw.WriteString(prefix + "<")
		sw.WriteString(xmlName)

		for k, v := range attrs { // write attrs
			sw.WriteString(fmt.Sprintf(" %s=\"%s\"", k, v))
		}

		sw.WriteString(">" + s.buildComment(withComment))
		prefix2 := prefix + "  "
		// sort keys
		var keys []string
		for k := range s.Properties {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		// build xml
		for _, k := range keys {
			v := s.Properties[k]
			if v.Name == "XMLName" && v.FullName == "xml.Name" { // ignore xmlname
				continue
			}
			key, _, isAttr, isOmitempty, isInner := v.XMLTag()
			if isAttr { // ignore attr
				continue
			}
			if key == "-" || (isOmitempty && v.ExampleValue() == NULL) { // ignore "-", omitempty
				continue
			}
			if isInner { // innerxml
				sw.WriteString(prefix2 + strings.Trim(v.ExampleValue(), "\""))
				sw.WriteString(v.buildComment(withComment, "innerxml") + "\n")
				continue
			}
			if len(v.Properties) == 0 && v.Type != ARRAY {
				example := v.ExampleValue()
				if example != NULL {
					// write xml node one line
					sw.WriteString(fmt.Sprintf("\n"+prefix2+"<%s>", key)) // write key
					sw.WriteString(strings.Trim(example, "\""))
					sw.WriteString(fmt.Sprintf("</%s>%s", key, v.buildComment(withComment))) // write key
				}
			} else {
				v.parseXML(depth+1, sw, true, withComment)
			}
		}

		sw.WriteString("\n" + prefix + fmt.Sprintf("</%s>", xmlName))
	} else if s.Type == ARRAY && s.ArraySchema != nil {
		if s.ArraySchema.Type == OBJECT {
			s.ArraySchema.parseXML(depth, sw, true, withComment)
		} else {
			sw.WriteString(fmt.Sprintf("\n%s<%s>%s\n", prefix, s.Name, s.buildComment(withComment))) // write key
			s.ArraySchema.parseXML(depth+1, sw, false, withComment)
			sw.WriteString(fmt.Sprintf("\n%s</%s>", prefix, s.Name)) // write key
		}
	} else { // write example value
		if isNewLine {
			sw.WriteString("\n")
		}
		example := s.ExampleValue()
		if example != NULL {
			sw.WriteString(prefix + "  " + strings.Trim(example, "\""))
		}
	}
}

func (v *TypeSchema) buildComment(withComment bool, withPrefix ...string) string {
	if v == nil {
		return ""
	}
	if !withComment {
		return ""
	}
	s := "  //"
	if len(withPrefix) > 0 {
		for _, x := range withPrefix {
			s += x + ", "
		}
	}
	if v.Type == ARRAY {
		arrayName := v.ArraySchema.Type // int
		if v.ArraySchema.Type == OBJECT {
			arrayName = v.ArraySchema.FullName
		}
		s += fmt.Sprintf("%s[%s]", ARRAY, arrayName)
	} else if v.Type == OBJECT { // object
		s += v.Type
		if v.FullName != "" {
			s += fmt.Sprintf("(%s)", v.FullName)
		}
	} else {
		s += v.Type
	}

	validate := v.ValidateTag()
	if validate != "" && !strings.Contains(validate, "required") {
		if v.IsRequired() {
			s += ", required"
		}
	}
	if validate != "" {
		s += fmt.Sprintf(", validate:\"%s\"", validate)
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

func (v *TypeSchema) hasJSONTag() bool {
	_, has := v.GetTag("json")
	return has
}

func (v *TypeSchema) ExampleValue() string {
	if v.example != "" {
		return v.example
	}
	if v.Type == ARRAY && v.ArraySchema == nil {
		return NULL
	}
	example := ""
	if val, has := v.GetTag("example"); has {
		example = val
	}
	if v.Type == OBJECT && (v.Properties == nil || len(v.Properties) == 0) {
		if v.FullName == "time.Time" {
			if example != "" {
				v.example = fmt.Sprintf("\"%s\"", example)
				return v.example
			}
			v.example = getTypeExample(v.FullName, example)
			return v.example
		}
		return NULL
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

func (v *TypeSchema) hasXMLName() (xmlName string, has bool) {
	if v.xmlName != "" {
		return v.xmlName, true
	}
	if v.Properties != nil {
		if x, ok := v.Properties["xmlname"]; ok {
			if x.Name == "XMLName" && x.FullName == "xml.Name" {
				val, has2 := x.GetTag("xml")
				if has2 {
					xmlName = strings.Split(val, ",")[0]
					has = true
					return
				}
			}
		}
	}
	return "", false
}

func (v *TypeSchema) XMLName() string {
	if v.xmlName != "" {
		return v.xmlName
	}
	if v.Properties != nil {
		if x, ok := v.Properties["xmlname"]; ok {
			if x.Name == "XMLName" && x.FullName == "xml.Name" {
				val, has := x.GetTag("xml")
				if has {
					return strings.Split(val, ",")[0]
				}
			}
		}
	}
	return v.Name
}

func (v *TypeSchema) XMLAttrs() map[string]string {
	if v.Properties != nil {
		attrsMap := make(map[string]string)
		for _, schema := range v.Properties {
			if schema.Name == "XMLName" && schema.FullName == "xml.Name" {
				continue
			}
			xmlTag, _, isAttr, _, _ := schema.XMLTag()
			if isAttr {
				attrsMap[xmlTag] = schema.ExampleValue()
			}
		}
		return attrsMap
	}
	return nil
}

func (v *TypeSchema) XMLTag() (xmlTag string, hasTag, isAttr, isOmitempty, isInner bool) {
	val, has := v.GetTag("xml")
	if has {
		arr := strings.Split(val, ",")
		if len(arr) > 0 {
			xmlTag = arr[0]
			hasTag = xmlTag != "" && xmlTag != "-"
			for i, a := range arr {
				if i == 0 {
					continue
				}

				switch a {
				case "attr":
					isAttr = true
				case "omitempty":
					isOmitempty = true
				case "innerxml":
					isInner = true
				}
			}
			return
		}
	}
	return v.Name, false, false, false, false
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
	case "time.Time":
		if example != "" {
			return fmt.Sprintf("\"%s\"", example)
		}
		t := time.Date(2022, 5, 16, 16, 47, 48, 741899000, time.Local) // use this time, prevent changes everytime build docs
		b, _ := t.MarshalJSON()
		return string(b)
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
		return "string"
	}
	return example
}

func getAllTagValue(field *ast.Field) string {
	if field != nil && field.Tag != nil && field.Tag.Value != "" {
		return strings.ReplaceAll(field.Tag.Value, "`", "")
	}
	return ""
}
