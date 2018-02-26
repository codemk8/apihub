build:
	go build -o bin/ap main.go

test:
	go test pkg/k8s
