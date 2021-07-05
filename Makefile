make:
	go run main.go

test:
	go test ./... -short

build:
	go build -o fm main.go

install:
	go install