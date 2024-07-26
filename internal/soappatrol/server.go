package soappatrol

import (
	"context"
	"encoding/xml"
	"github.com/dottorblaster/soappatrol/pkg/soap"
	"net"
	"net/http"
	"time"

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

type Server struct {
	SoapServer *http.Server
	Logger     *zap.SugaredLogger
}

func New(config Config, logger *zap.SugaredLogger) Server {
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

	httpServer := &http.Server{
		Handler:           soapServer,
		ReadHeaderTimeout: 32 * time.Second,
	}

	soappatrolServer := Server{
		SoapServer: httpServer,
		Logger:     logger,
	}

	return soappatrolServer
}

func (s *Server) ListenAndServe(socket string) error {
	unixListener, err := net.Listen("unix", socket)
	if err != nil {
		s.Logger.Errorw("Error listening on the socket", err)
		return err
	}

	err = s.SoapServer.Serve(unixListener)
	if err != http.ErrServerClosed {
		s.Logger.Errorw("Error serving on the listener")
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	err := s.SoapServer.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}
