build:
	go build -o dag-go main.go

run:
	go build -o dag-go main.go && ./dag-go

test:
	go test -v ./...