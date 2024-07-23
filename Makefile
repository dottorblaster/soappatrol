soappatrol:
	CGO_ENABLED=0 go build -o soappatrol cmd/soappatrol/main.go

.PHONY: format
format:
	go fmt ./...

all: soappatrol

.PHONY: test
test:
	go test -v -p 1 -race ./...

clean:
	rm -rf soappatrol
