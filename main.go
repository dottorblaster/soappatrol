package main

import (
	"encoding/xml"
	"fmt"
	"github.com/dottorblaster/soappatrol/soap"
	"net"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

// FooRequest a simple request
type FooRequest struct {
	XMLName xml.Name `xml:"fooRequest"`
	Foo     string
}

type MockRequest struct {
  XMLName xml.Name
}

// FooResponse a simple response
type FooResponse struct {
	Bar string
}

type Request struct {
	Action   string
	Tagname  string
	Response string
}

type Config struct {
	Requests []Request `toml:"request"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock config.toml")
		return
	}

	configString, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	if _, err := toml.Decode(string(configString), &config); err != nil {
		fmt.Println(os.Stderr, "Unable to parse config")
	}

	fmt.Println("Unix SOAP server")

	os.Remove(os.Args[1])

	soapServer := soap.NewServer()

	for _, r := range config.Requests {
		fmt.Printf("Registering: %s %s\n", r.Action, r.Tagname)

		soapServer.RegisterHandler(
			"/",
			// SOAPAction
			r.Action,
			// tagname of soap body content
			r.Tagname,
			// RequestFactoryFunc - give the server sth. to unmarshal the request into
			func() interface{} {
				return &MockRequest{}
			},
			// OperationHandlerFunc - do something
			func(request interface{},
				w http.ResponseWriter,
				httpRequest *http.Request,
			) (response interface{}, err error) {
				response = r.Response
				return
			},
		)

	}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		panic(err)
	}

	http.Serve(unixListener, soapServer)
}
