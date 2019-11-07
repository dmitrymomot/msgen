/*
Copyright © 2019 Dmytro Momot

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	_ "github.com/dmitrymomot/msgen/statik" // import templates
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n", "init", "create"},
	Short:   "Generate a new microservice from template",
	Example: "msgen new --grpc --grpc_port=9876 --pub --path=srv/users github.com/your-username/users-srv",
	Run:     run,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a service name")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("namespace", "n", "default", "service namespace")
	newCmd.Flags().String("path", "", "create in directory")

	newCmd.Flags().Bool("grpc", false, "set up gRPC server")
	newCmd.Flags().Int("grpc_port", 9200, "grpc port")

	newCmd.Flags().Bool("twirp", false, "set up twirp server (exposed via http_port)")

	newCmd.Flags().Bool("http", false, "set up http server")
	newCmd.Flags().Int("http_port", 8888, "http port")

	newCmd.Flags().Bool("pub", false, "set up NATS publisher")
	newCmd.Flags().Bool("sub", false, "set up NATS subscriber")

	newCmd.Flags().Bool("redis_pool", false, "set up redis connections pool")
	newCmd.Flags().String("redis_host", "", "redis connection path")
	// newCmd.Flags().Bool("jobs", false, "set up jobs queue based on redis")
	newCmd.Flags().StringSlice("jobs", nil, "set up jobs queue based on redis (e.g.: --jobs=send_email,send_notification)")
}

type (
	options struct {
		ServiceName string
		ServicePath string
		Namespace   string
		Path        string
		AbsPath     string

		Linkerd bool

		Grpc        bool
		GrpcLB      bool
		GrpcPort    int
		GrpcClients []string
		GrpcSvc     bool

		Twirp        bool
		TwirpClients []string

		RPC        bool
		RPCMethods []string

		HTTP          bool
		HTTPPort      int
		HTTPEndpoints []string
		HTTPHealth    string
		HTTPLB        bool
		HTTPSvc       bool

		RedisPool bool
		RedisHost string
		Jobs      []string

		DB *dbOptions

		ClientHelper bool
		K8s          bool

		TLS     bool
		TLSPath string

		Nats bool
		Pub  bool
		Sub  bool

		CustomOptions map[string]interface{}
	}

	file struct {
		Template string
		Path     string
		Options  map[string]interface{}
	}

	dbOptions struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
	}
)

func newFile(t string, opt options) file {
	f := file{Template: t, Path: t}
	if strings.Contains(t, "examplesrv") {
		f.Path = strings.ReplaceAll(t, "examplesrv", PackageName(opt.ServiceName))
	}
	return f
}

func run(cmd *cobra.Command, args []string) {
	sfs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	opt := options{}
	opt.ServicePath = cmd.Flags().Arg(0)
	opt.ServiceName = filepath.Base(opt.ServicePath)
	opt.Namespace, _ = cmd.Flags().GetString("namespace")
	opt.Grpc, _ = cmd.Flags().GetBool("grpc")
	opt.GrpcPort, _ = cmd.Flags().GetInt("grpc_port")
	opt.Twirp, _ = cmd.Flags().GetBool("twirp")
	opt.HTTP, _ = cmd.Flags().GetBool("http")
	opt.HTTPPort, _ = cmd.Flags().GetInt("http_port")
	opt.Jobs, _ = cmd.Flags().GetStringSlice("jobs")

	log.Println(opt.Jobs)

	opt.Path, _ = cmd.Flags().GetString("path")
	opt.Path = strings.Trim(opt.Path, "/")
	if opt.Path == "" || opt.Path == "." {
		opt.Path = opt.ServiceName
	}
	opt.AbsPath, _ = filepath.Abs(opt.Path)

	files := []file{
		// newFile(".gitignore", opt),
		newFile("Dockerfile", opt),
		newFile("go.mod", opt),
		newFile("main.go", opt),
		newFile("Makefile", opt),
		newFile("README.md", opt),
		newFile("pb/examplesrv/service.proto", opt),
		newFile("service/service.go", opt),
	}

	for _, file := range files {
		if err := execFile(sfs, file, opt); err != nil {
			log.Fatalf("file: %+v, options: %+v, error: %v", file, opt, err)
		}
	}

	if opt.Jobs != nil && len(opt.Jobs) > 0 {
		if err := execFile(sfs, newFile("jobs/worker.go", opt), opt); err != nil {
			panic(err)
		}
		for _, job := range opt.Jobs {
			filename := strings.ReplaceAll(filepath.Base(strings.ToLower(job)), "-", "_") + ".go"
			fstr := file{
				Template: "jobs/job.go",
				Path:     filepath.Join("jobs", filename),
				Options: map[string]interface{}{
					"jobTitle": ToTitle(job),
					"jobName":  job,
				},
			}
			if err := execFile(sfs, fstr, opt); err != nil {
				panic(err)
			}
		}
	}

	if opt.RPCMethods != nil && len(opt.RPCMethods) > 0 {
		for _, method := range opt.RPCMethods {
			filename := strings.ReplaceAll(filepath.Base(strings.ToLower(method)), "-", "_") + ".go"
			fstr := file{
				Template: "service/method.go",
				Path:     filepath.Join("service", filename),
				Options: map[string]interface{}{
					"methodTitle": ToTitle(method),
					"methodName":  method,
				},
			}
			if err := execFile(sfs, fstr, opt); err != nil {
				panic(err)
			}
		}
	}

	if opt.ClientHelper {
		if err := execFile(sfs, newFile("client/client.go", opt), opt); err != nil {
			panic(err)
		}
	}

	if opt.DB != nil {
		if err := execFile(sfs, newFile("db_logger.go", opt), opt); err != nil {
			panic(err)
		}
	}

	if opt.K8s {
		if err := execFile(sfs, newFile("k8s.yml", opt), opt); err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		if err := execFile(sfs, file, opt); err != nil {
			panic(err)
		}
	}
}

func execFile(sfs http.FileSystem, f file, opt options) error {
	b, err := fs.ReadFile(sfs, "/"+f.Template)
	if err != nil {
		return err
	}

	t, err := template.New("f").Funcs(template.FuncMap{
		"title":   ToTitle,
		"camel":   ToCamelCase,
		"package": PackageName,
		"url2var": URLToVarName,
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

func prepareString(str string) string {
	str = strings.ReplaceAll(str, "_", " ")
	str = strings.ReplaceAll(str, "-", " ")
	str = strings.ReplaceAll(str, "/", " ")
	str = strings.ReplaceAll(str, ".", " ")
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
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	return ToCamelCase(prepareString(u.Hostname()))
}
