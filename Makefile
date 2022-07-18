build:
	echo "Build agent"
	go build -o bin/agent cmd/agent/main.go
	echo "Build agent"
	go build -o bin/server cmd/server/main.go

server:
	go run cmd/server/main.go

agent:
	go run cmd/agent/main.go

run: server
