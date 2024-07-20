package soappatrol

import (
	"encoding/xml"
	"github.com/dottorblaster/soappatrol/pkg/soap"
	"net/http"

	"go.uber.org/zap"
)

type Config struct {
	Requests []Request `toml:"request"`
}

type Request struct {
	Action   string
	Tagname  string
	Response string
}

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

type MockResponse struct {
	Response string `xml:",innerxml"`
}

func New(config Config, logger *zap.SugaredLogger) *soap.Server {
	soapServer := soap.NewServer()

	for _, r := range config.Requests {
		logger.Infow("Registering: ", r.Action, r.Tagname)

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

	return soapServer
}
