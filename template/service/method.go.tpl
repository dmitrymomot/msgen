package service

import (
	"context"
	"{{ .ServicePath }}/pb/{{ package .ServiceName }}"
)

func (s *service) {{ index .CustomOptions "methodTitle" }}(ctx context.Context, req *{{ package .ServiceName }}.{{ index .CustomOptions "methodTitle" }}Req) (*{{ package .ServiceName }}.{{ index .CustomOptions "methodTitle" }}Resp, error) {
	s.log.Info().Str("str", req.GetStr()).Msg("received string")
	return &{{ package .ServiceName }}.{{ index .CustomOptions "methodTitle" }}Resp{Str: "{{ package .ServiceName }}.{{ index .CustomOptions "methodTitle" }}: received: " + req.GetStr()}, nil
}
