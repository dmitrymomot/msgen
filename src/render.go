package src

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

func render(sfs http.FileSystem, f file, opt options) error {
	b, err := fs.ReadFile(sfs, "/"+f.Template)
	if err != nil {
		return err
	}

	t, err := template.New("f").Funcs(template.FuncMap{
		"title":    ToTitle,
		"camel":    ToCamelCase,
		"package":  PackageName,
		"urlTovar": URLToVarName,
		"asis": func(s string) template.HTML {
			return template.HTML(s)
		},
	}).Parse(string(b))
	if err != nil {
		return errors.Wrap(err, "parse template")
	}

	absFilePath := filepath.Join(opt.AbsPath, f.Path)
	absFileDir := filepath.Dir(absFilePath)

	// check if dir exists
	if _, err := os.Stat(absFileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(absFileDir, os.ModePerm|os.ModeDir); err != nil {
			panic(err)
		}
	}

	fmt.Println("Creating of file " + absFilePath)
	fl, err := os.Create(absFilePath)
	if err != nil {
		return errors.Wrap(err, "create file")
	}

	defer fl.Close()

	if f.Options != nil {
		opt.CustomOptions = f.Options
	}

	if strings.Contains(f.Path, ".go") {
		var out bytes.Buffer
		err = t.Execute(&out, opt)
		if err != nil {
			return errors.Wrapf(err, "Could not process template %s", f.Path)
		}

		b, err := format.Source(out.Bytes())
		if err != nil {
			fmt.Print(string(out.Bytes()))
			return errors.Wrapf(err, "\nCould not format Go file %s\n", f.Path)
		}

		if _, err = fl.Write(b); err != nil {
			return err
		}
	} else {
		if err := t.Execute(fl, opt); err != nil {
			return errors.Wrap(err, "execute template")
		}
	}

	return nil
}
