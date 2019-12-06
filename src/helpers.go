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
	return goNamingConvention(str)
}

// PackageName ...
func PackageName(str string) string {
	return strings.ReplaceAll(strings.ToLower(prepareString(str)), " ", "")
}

// GRPCPackageName ...
func GRPCPackageName(str string) string {
	return strings.ReplaceAll(strings.ToLower(prepareString(str)), " ", ".")
}

// ToURL ...
func ToURL(str string) string {
	return strings.ReplaceAll(strings.ToLower(prepareString(str)), " ", "-")
}

// ToCamelCase ...
func ToCamelCase(str string) string {
	return goNamingConvention(LcFirst(ToTitle(str)))
}

// UcFirst ...
func UcFirst(str string) string {
	for i, v := range str {
		return goNamingConvention(string(unicode.ToUpper(v)) + str[i+1:])
	}
	return ""
}

// LcFirst ...
func LcFirst(str string) string {
	for i, v := range str {
		return goNamingConvention(string(unicode.ToLower(v)) + str[i+1:])
	}
	return ""
}

// URLToVarName ...
func URLToVarName(str string) string {
	str = alphaStr.ReplaceAllString(str, " ")
	return goNamingConvention(ToCamelCase(prepareString(str)))
}

func goNamingConvention(str string) string {
	str = strings.ReplaceAll(str, "Api", "API")
	str = strings.ReplaceAll(str, "Rpc", "RPC")
	str = strings.ReplaceAll(str, "Http", "HTTP")
	str = strings.ReplaceAll(str, "Id", "ID")
	str = strings.ReplaceAll(str, "Url", "URL")
	str = strings.ReplaceAll(str, "Uri", "URI")
	return str
}
