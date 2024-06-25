soappatrol:
	CGO_ENABLED=0 go build -o soappatrol main.go

client:
	CGO_ENABLED=0 go build -o client client.go

all: soappatrol client

clean:
	rm -rf soappatrol
	rm -rf client
