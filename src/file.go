package src

import (
	"strings"
)

type file struct {
	Template string
	Path     string
	Options  map[string]interface{}
}

func newFile(templatePath, serviceName string) file {
	f := file{Template: templatePath, Path: strings.TrimSuffix(templatePath, ".tpl")}
	if strings.Contains(f.Path, "examplesrv") {
		f.Path = strings.ReplaceAll(f.Path, "examplesrv", PackageName(serviceName))
	}
	return f
}

func newCustomFile(templatePath, filename, serviceName string, opt map[string]interface{}) file {
	f := file{
		Template: templatePath,
		Path:     filename,
		Options:  opt,
	}
	if strings.Contains(f.Path, "examplesrv") {
		f.Path = strings.ReplaceAll(f.Path, "examplesrv", PackageName(serviceName))
	}
	return f
}

func getDefaultFilesList(serviceName string) []file {
	return []file{
		file{Template: "gitignore.tpl", Path: ".gitignore"},
		newFile("Dockerfile.tpl", serviceName),
		newFile("go.mod.tpl", serviceName),
		newFile("logger.go.tpl", serviceName),
		newFile("main.go.tpl", serviceName),
		newFile("Makefile.tpl", serviceName),
		newFile("README.md.tpl", serviceName),
	}
}
