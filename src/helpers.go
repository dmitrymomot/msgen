package src

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	alphaStr = regexp.MustCompile("[^a-zA-Z]+")
)

func prepareString(str string) string {
	str = alphaStr.ReplaceAllString(str, " ")
	return str
}

// ToTitle ...
func ToTitle(str string) string {
	str = strings.Title(prepareString(str))
	str = strings.ReplaceAll(str, " ", "")
	return str
}

// PackageName ...
func PackageName(str string) string {
	return strings.ReplaceAll(prepareString(str), " ", "")
}

// ToCamelCase ...
func ToCamelCase(str string) string {
	return LcFirst(ToTitle(str))
}

// UcFirst ...
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LcFirst ...
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// URLToVarName ...
func URLToVarName(str string) string {
	str = alphaStr.ReplaceAllString(str, " ")
	return ToCamelCase(prepareString(str))
}
