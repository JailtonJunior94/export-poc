build:
	@echo "Compiling..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin/export ./cmd/main.go