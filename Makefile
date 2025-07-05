.PHONY : format clean install build

run:
	go run ./cmd/webapi/main.go

format:
	gofmt -s -w .

clean:
	go mod tidy

install:
	go mod download

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
