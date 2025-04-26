all:
	go build -o simple-formater
fmt:
	@find . -name "*.go" -exec go fmt {} \;
