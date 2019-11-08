package src

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/dmitrymomot/msgen/statik" // import templates
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

// Generate project
func Generate(cmd *cobra.Command, args []string) {
	sfs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	opt := parseOptions(cmd, args)
	files := getDefaultFilesList(opt.ServiceName)

	if opt.RPC {
		files = append(files, newFile("pb/examplesrv/service.proto.tpl", opt.ServiceName))
		files = append(files, newFile("service/service.go.tpl", opt.ServiceName))
	}

	if opt.RedisPool {
		files = append(files, newFile("redis.go.tpl", opt.ServiceName))
	}

	if opt.Jobs != nil && len(opt.Jobs) > 0 {
		files = append(files, newFile("jobs/worker.go.tpl", opt.ServiceName))
		for _, job := range opt.Jobs {
			filename := strings.ReplaceAll(filepath.Base(strings.ToLower(job)), "-", "_") + ".go"
			files = append(files, newCustomFile("jobs/job.go.tpl", filename, opt.ServiceName, map[string]interface{}{
				"jobTitle": ToTitle(job),
				"jobName":  job,
			}))
		}
	}

	if opt.RPCMethods != nil && len(opt.RPCMethods) > 0 {
		for _, method := range opt.RPCMethods {
			filename := filepath.Join("service", strings.ReplaceAll(filepath.Base(strings.ToLower(method)), "-", "_")+".go")
			files = append(files, newCustomFile("service/method.go.tpl", filename, opt.ServiceName, map[string]interface{}{
				"methodTitle": ToTitle(method),
				"methodName":  method,
			}))
		}
	}

	if opt.ClientHelper {
		files = append(files, newFile("client/client.go.tpl", opt.ServiceName))
	}

	if opt.DB != nil {
		files = append(files, newFile("db_logger.go.tpl", opt.ServiceName))
	}

	if opt.K8s {
		files = append(files, newFile("k8s.yml.tpl", opt.ServiceName))
	}

	if opt.Nats {
		files = append(files, newFile("nats.go.tpl", opt.ServiceName))
	}

	if err := handleFiles(sfs, files, opt); err != nil {
		log.Fatal(err)
	}
}

func handleFiles(sfs http.FileSystem, files []file, opt options) error {
	for _, file := range files {
		if err := render(sfs, file, opt); err != nil {
			return fmt.Errorf("file: %+v, options: %+v, error: %v", file, opt, err)
		}
	}
	return nil
}

func render(sfs http.FileSystem, f file, opt options) error {
	b, err := fs.ReadFile(sfs, "/"+f.Template)
	if err != nil {
		return err
	}

	t, err := template.New("f").Funcs(template.FuncMap{
		"title":   ToTitle,
		"camel":   ToCamelCase,
		"package": PackageName,
		"url2var": URLToVarName,
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
