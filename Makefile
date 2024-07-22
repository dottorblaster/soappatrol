soappatrol:
	CGO_ENABLED=0 go build -o soappatrol cmd/soappatrol/main.go

.PHONY: format
format:
	go fmt ./...

all: soappatrol

clean:
	rm -rf soappatrol
