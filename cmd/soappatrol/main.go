package main

import (
	"encoding/xml"
	"github.com/dottorblaster/soappatrol/pkg/soap"
	"net"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
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

type MockResponse struct {
	Response string `xml:",innerxml"`
}

type Config struct {
	Requests []Request `toml:"request"`
}

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic("Error initializing logger")
		}
	}()
	if len(os.Args) < 3 {
		logger.Errorw("usage:", os.Args[0], "/path.sock config.toml")
		return
	}

	configString, err := os.ReadFile(os.Args[2])
	if err != nil {
		logger.Error(zap.Error(err))
		return
	}

	var config Config
	if _, err := toml.Decode(string(configString), &config); err != nil {
		logger.Errorw("Unable to parse config")
	}

	logger.Infow("Unix SOAP server")

	os.Remove(os.Args[1])

	soapServer := soap.NewServer()

	for _, r := range config.Requests {
		logger.Infow("Registering: %s %s\n", r.Action, r.Tagname)

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
			func(_ interface{},
				_ http.ResponseWriter,
				_ *http.Request,
			) (interface{}, error) {
				response := &MockResponse{
					Response: r.Response,
				}
				return response, nil
			},
		)

	}

	unixListener, err := net.Listen("unix", os.Args[1])
	if err != nil {
		logger.Errorw("Error listening on the port")
		panic(err)
	}

	// We have to bypass http.Server here because we have to explicitly
	// bind our baked implementation of the SOAP server to the unix socket
	// nolint:gosec
	err = http.Serve(unixListener, soapServer)
	if err != nil {
		logger.Errorw("Error serving on the listener")
		panic(err)
	}
}
