build:
	go build -o bin/ap main.go

test:
	go test github.com/codemk8/apihub/pkg/k8s
	go test github.com/codemk8/apihub/pkg/helm
