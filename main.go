package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
        "github.com/globusdigital/soap"
)

// FooRequest a simple request
type FooRequest struct {
	XMLName xml.Name `xml:"fooRequest"`
	Foo     string
}

// FooResponse a simple response
type FooResponse struct {
	Bar string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock [wwwroot]")
		return
	}

	fmt.Println("Unix HTTP SOAP server")

	root := "."
	if len(os.Args) > 2 {
		root = os.Args[2]
	}

	os.Remove(os.Args[1])

        soapServer := soap.NewServer()

	soapServer.HandleOperation(
		// SOAPAction
		"operationFoo",
		// tagname of soap body content
		"fooRequest",
		// RequestFactoryFunc - give the server sth. to unmarshal the request into
		func() interface{} {
			return &FooRequest{}
		},
		// OperationHandlerFunc - do something
		func(request interface{}, w http.ResponseWriter, httpRequest *http.Request) (response interface{}, err error) {
			fooRequest := request.(*FooRequest)
			fooResponse := &FooResponse{
				Bar: "Hello " + fooRequest.Foo,
			}
			response = fooResponse
			return
		},
	)

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}

        soapServer.Serve(unixListener)
}
