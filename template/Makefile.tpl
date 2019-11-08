
# Environment variables
LATEST_COMMIT := $$(git rev-parse HEAD)


.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
{{if .RPC}}
.PHONY: proto
proto: ## Compile protobuf for golang
	@protoc -I /usr/local/include -I . \
		-I$(GOPATH)/src \
		-I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
		--go_out=plugins=grpc:. \
		--twirp_out=. \
		--validate_out=lang=go:. \
		pb/**/*.proto
{{end}}
.PHONY: build
build: {{if .RPC}}proto{{end}} ## Build application{{if .RPC}} and compile protobuf for golang{{end}}
	@go clean
	@CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-a -installsuffix nocgo \
	-ldflags "-X main.buildTag=`date -u +%Y%m%d.%H%M%S`-$(LATEST_COMMIT)" \
	-o {{ .ServiceName }} .

.PHONY: docker
docker: ## Build docker image
	@docker build . -t {{ .ServiceName }}:latest
{{if .K8s}}
.PHONY: deploy
deploy: ## Deploy pods to kubernetes
	@kubectl apply -f k8s.yml

.PHONY: down
down: ## Down pods
	@kubectl delete -f k8s.yml

.PHONY: reload
reload: down deploy info ## Reload after app was rebuilt

.PHONY: info
info: ## Get cluster info
	@kubectl get all

.PHONY: logs
log: ## Show logs
	@kubectl logs -lapp={{ .ServiceName }} --container={{ .ServiceName }}
{{end}}