package src

import (
	"path/filepath"
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
		newFile("main.go.tpl", serviceName),
		newFile("Makefile.tpl", serviceName),
		newFile("README.md.tpl", serviceName),
	}
}

func getFilesList(opt options) []file {
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

	return files
}
