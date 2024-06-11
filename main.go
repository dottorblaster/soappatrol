package main

import (
	"encoding/xml"
	"fmt"
	"github.com/foomo/soap"
	"net"
	"net/http"
	"os"
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
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock")
		return
	}

	fmt.Println("Unix SOAP server")

	os.Remove(os.Args[1])

	soapServer := soap.NewServer()

	soapServer.RegisterHandler(
		"/",
		// SOAPAction
		"operationFoo",
		// tagname of soap body content
		"fooRequest",
		// RequestFactoryFunc - give the server sth. to unmarshal the request into
		func() interface{} {
			return &FooRequest{}
		},
		// OperationHandlerFunc - do something
		func(request interface{},
			w http.ResponseWriter,
			httpRequest *http.Request,
		) (response interface{}, err error) {
			fooRequest := request.(*FooRequest)
			fooResponse := &FooResponse{
				Bar: "Hellow " + fooRequest.Foo,
			}
			response = fooResponse
			return
		},
	)

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}

	http.Serve(unixListener, soapServer)
}
