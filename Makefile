hello:
	echo "Hello"

build:
	go build -o bin/cicd-server cmd/webhookcicd/main.go

run:
	SECRET=my-webhook go run cmd/webhookcicd/main.go

proto:
	protoc --proto_path=internal/proto  --go_out=plugins=grpc:internal client.proto

distribute:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/cicd-server-linux-arm cmd/webhookcicd/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/cicd-server-linux-arm64 cmd/webhookcicd/main.go

all: build