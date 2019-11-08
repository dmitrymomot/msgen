kind: ServiceAccount
apiVersion: v1
metadata:
  name: {{ .ServiceName }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: {{ .ServiceName }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .ServiceName }}
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      {{if .Linkerd}}annotations:
        config.linkerd.io/skip-outbound-ports: "4222,5432,6379"
        linkerd.io/inject: enabled{{end}}
      creationTimestamp: null
      labels:
        app: {{ .ServiceName }}
    spec:
      serviceAccountName: {{ .ServiceName }}
      containers:
        - name: {{ .ServiceName }}
          image: {{ .ServiceName }}:latest
          imagePullPolicy: IfNotPresent
          args: [
            "-debug=true",
            "-app_name='The Awesome App'",
            {{if .HTTP}}"-http_port={{.HTTPPort}}",{{end}}
            {{if .Grpc}}"-grpc_port={{.GrpcPort}}",{{end}}
            {{if .DB}}"-db_host={{.DB.Host}}",
            "-db_port={{.DB.Port}}",
            "-db_name={{.DB.Name}}",
            "-db_user={{.DB.User}}",
            "-db_password={{.DB.Password}}",
            "-db_pool_size=10",{{end}}
            {{if .RedisPool}}"-redis_host={{.RedisHost}}",{{end}}
            {{if .Nats}}"-nats_host=nats://nats-cluster:4222",
            "-nats_queue_subject={{ .ServiceName }}"{{end}}
          ]
          ports:
            {{if .Grpc}}- containerPort: {{.GrpcPort}}
              name: grpc{{end}}
            {{if .HTTP}}- containerPort: {{.HTTPPort}}
              name: http{{end}}
status: {}
{{if .GrpcSrv}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .ServiceName }}
spec:
  selector:
    app: {{ .ServiceName }}
  clusterIP: None
  ports:
    - name: grpc
      port: {{.GrpcPort}}
      targetPort: {{.GrpcPort}}
{{end}}
{{if .GrpcLB}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .ServiceName }}-grpc-lb
spec:
  type: LoadBalancer
  selector:
    app: {{ .ServiceName }}
  ports:
    - name: grpc
      port: {{.GrpcPort}}
      targetPort: {{.GrpcPort}}
{{end}}
{{if .HTTPSrv}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .ServiceName }}
spec:
  selector:
    app: {{ .ServiceName }}
  clusterIP: None
  ports:
    - name: http
      port: {{.HTTPPort}}
      targetPort: {{.HTTPPort}}
{{end}}
{{if .HTTPLB}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .ServiceName }}-http-lb
spec:
  type: LoadBalancer
  selector:
    app: {{ .ServiceName }}
  ports:
    - name: http
      port: {{.HTTPPort}}
      targetPort: {{.HTTPPort}}
{{end}}