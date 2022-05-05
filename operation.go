package apidoc

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
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
	case descriptionAttr:
		operation.ParseDescriptionComment(lineRemainder)
	case tagsAttr:
		operation.ParseTagsComment(lineRemainder)
	case apiAttr:
		return operation.ParseRouterComment(lineRemainder)
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
	return nil
}
