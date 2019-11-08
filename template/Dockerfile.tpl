FROM alpine:latest
RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/* /tmp/*
ADD {{ .ServiceName }} /{{ .ServiceName }}
ENTRYPOINT [ "/{{ .ServiceName }}" ]
