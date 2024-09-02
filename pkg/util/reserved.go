package util

import "strings"

var reservedWords = []string{
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"goto",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",
}

var basicTypes = []string{
	"bool",
	"string",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"byte",
	"rune",
	"float32",
	"float64",
	"complex64",
	"complex128",
}

var basicDataStructures = []string{
	"map",
	"[]",
}

func IsReservedWord(word string) bool {
	word = strings.ToLower(word)
	for _, reserved := range reservedWords {
		if reserved == word {
			return true
		}
	}
	return false
}

func IsBasicType(variableType string) bool {
	for _, basicType := range basicTypes {
		if basicType == variableType {
			return true
		}
	}
	return false
}

func IsBasicDataStructure(variableType string) bool {
	for _, basicDataStructure := range basicDataStructures {
		if strings.HasPrefix(variableType, basicDataStructure) {
			return true
		}
	}
	return false
}
