package parser

import (
	"fmt"
)

const (
	// ARRAY represent a array value.
	ARRAY = "array"
	// OBJECT represent a object value.
	OBJECT = "object"
	// BOOLEAN represent a boolean value.
	BOOLEAN = "boolean"
	// INTEGER represent a integer value.
	INTEGER = "integer"
	// NUMBER represent a number value.
	NUMBER = "number"
	// STRING represent a string value.
	STRING = "string"
	// FUNC represent a function value.
	FUNC = "func"
	// INTERFACE represent a interface value.
	INTERFACE = "interface{}"
	// ANY represent a any value.
	ANY = "any"
	// NIL represent a empty value.
	NIL  = "nil"
	NULL = "null"
)

// CheckSchemaType checks if typeName is not a name of primitive type.
func CheckSchemaType(typeName string) error {
	if !IsPrimitiveType(typeName) {
		return fmt.Errorf("%s is not basic types", typeName)
	}

	return nil
}

// IsSimplePrimitiveType determine whether the type name is a simple primitive type.
func IsSimplePrimitiveType(typeName string) bool {
	switch typeName {
	case STRING, NUMBER, INTEGER, BOOLEAN:
		return true
	}

	return false
}

// IsPrimitiveType determine whether the type name is a primitive type.
func IsPrimitiveType(typeName string) bool {
	switch typeName {
	case STRING, NUMBER, INTEGER, BOOLEAN, ARRAY, OBJECT, FUNC:
		return true
	}

	return false
}

// IsNumericType determines whether the type name is a numeric type.
func IsNumericType(typeName string) bool {
	return typeName == INTEGER || typeName == NUMBER
}

// TransToValidSchemeType indicates type will transfer golang basic type to supported type.
func TransToValidSchemeType(typeName string) string {
	switch typeName {
	case "uint", "int", "uint8", "int8", "uint16", "int16", "byte":
		return INTEGER
	case "uint32", "int32", "rune":
		return INTEGER
	case "uint64", "int64":
		return INTEGER
	case "float32", "float64":
		return NUMBER
	case "bool":
		return BOOLEAN
	case "string":
		return STRING
	}

	return typeName
}

// IsGolangPrimitiveType determine whether the type name is a golang primitive type.
func IsGolangPrimitiveType(typeName string) bool {
	switch typeName {
	case "uint",
		"int",
		"uint8",
		"int8",
		"uint16",
		"int16",
		"byte",
		"uint32",
		"int32",
		"rune",
		"uint64",
		"int64",
		"float32",
		"float64",
		"bool",
		"string",
		"any":
		return true
	}

	return false
}
