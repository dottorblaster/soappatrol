soappatrol:
	go build -o soappatrol main.go

client:
	go build -o client client.go

all: soappatrol client

clean:
	rm -rf soappatrol
	rm -rf client
