soappatrol:
	CGO_ENABLED=0 go build -o soappatrol cmd/soappatrol/main.go

all: soappatrol

clean:
	rm -rf soappatrol
