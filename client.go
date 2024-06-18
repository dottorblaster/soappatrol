package main

import (
	"context"
	"encoding/xml"
	"net"
        "time"
	"net/http"
	"os"
        "fmt"

	"github.com/hooklift/gowsdl/soap"
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

func main() {
	socket := os.Args[1]
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
          fmt.Println(err)
	}

	fmt.Println(response)
}
