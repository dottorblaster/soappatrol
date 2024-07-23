package soappatrol_test

import (
	"context"
	"encoding/xml"
	"github.com/dottorblaster/soappatrol/internal/soappatrol"
	"github.com/hooklift/gowsdl/soap"
	"go.uber.org/zap/zaptest"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
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

type GetInstanceProperties struct {
	XMLName xml.Name `xml:"urn:SAPControl GetInstanceProperties"`
}

type GetInstancePropertiesResponse struct {
	XMLName    xml.Name            `xml:"urn:SAPControl GetInstancePropertiesResponse"`
	Properties []*InstanceProperty `xml:"properties>item,omitempty" json:"properties>item,omitempty"`
}

type InstanceProperty struct {
	Property     string `xml:"property,omitempty" json:"property,omitempty"`
	Propertytype string `xml:"propertytype,omitempty" json:"propertytype,omitempty"`
	Value        string `xml:"value,omitempty" json:"value,omitempty"`
}

func TestSoappatrolServer(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	defer func() {
		err := logger.Sync()
		if err != nil {
			panic("Error initializing logger")
		}
	}()

	socket := "test_socket_1234"

	config := soappatrol.Config{Requests: []soappatrol.Request{
		{
			Action:  "''",
			Tagname: "GetInstanceProperties",
			Response: `
                          <GetInstancePropertiesResponse xmlns="urn:SAPControl">
                            <properties>
                              <item>
                                <property>Pasta</property>
                                <propertytype>Meal</propertytype>
                                <value>Carbonara</value>
                              </item>
                            </properties>
                          </GetInstancePropertiesResponse>
          `,
		},
	},
	}

	soapServer := soappatrol.New(config, logger)

	go func(server soappatrol.Server, socket string) {
		err := server.ListenAndServe(socket)
		if err != nil {
			t.Errorf("Error during server listen")

		}
	}(soapServer, socket)

	time.Sleep(3 * time.Second)

	udsClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", socket)

			},
		},
	}

	// The url used here is just phony:
	// we need a well formed url to create the instance but the above DialContext function won't actually use it.
	client := soap.NewClient("http://unix", soap.WithHTTPClient(udsClient))

	request := &GetInstanceProperties{}
	response := &GetInstancePropertiesResponse{}

	err := client.Call("''", request, response)
	if err != nil {
		t.Errorf("Error during client call")
	}

	if response.Properties[0].Property != "Pasta" {
		t.Errorf("megafail")
	}

	os.Remove(socket)
}
