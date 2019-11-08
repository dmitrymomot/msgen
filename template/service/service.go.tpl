package service

import (
	"{{ .ServicePath }}/pb/{{ package .ServiceName }}"
	{{ if .DB }}"github.com/go-pg/pg/v9"{{ end }}
	{{ if .Jobs }}"github.com/gocraft/work"{{ end }}
)

type service struct {
	{{ package .ServiceName }}.UnimplementedServiceServer
	log *Logger // uber's structured and leveled logger
	{{ if .DB }}db  *pg.DB             // postgresql database connection{{ end }}
	{{ if .Jobs }}wp  *work.Enqueuer     // redis worker pool{{ end }}
}

// New service factory
func New(log *Logger{{ if .DB }}, db *pg.DB{{ end }}{{ if .Jobs }}, wp *work.Enqueuer{{ end }}) {{ package .ServiceName }}.ServiceServer {
	return &service{log: log{{ if .DB }}, db: db{{ end }}{{ if .Jobs }}, wp: wp{{ end }}}
}
