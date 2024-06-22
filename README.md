# soappatrol

**soappatrol** is a program written in the Go programming language that spawns a SOAP mock server generated from a configuration file. It allows you to simulate a SOAP server for testing and development purposes, having it listening on a Unix socket of your choice.

## Features

- Generates a SOAP mock server from a configuration file
- Supports multiple SOAP endpoints and responses
- Listens on a UNIX socket
- Easy to use command-line interface

## Installation

To install soappatrol, you will need to have Go installed on your system. You can then install it using the Go package manager:
Text Only

go install github.com/dottorblaster/soappatrol@latest

## Usage

To use soappatrol, you need to create a configuration file that defines the SOAP endpoints and responses. Here's an example configuration file:

```toml
[[request]]
action = "operationFoo"
tagname = "barRequest"
response = """
<HelloResponse xmlns="http://example.com/">
  <HelloResult>Hello, world!</HelloResult>
</HelloResponse>
"""
```

To start the SOAP mock server, run the following command:

```
soappatrol /path/to/socket.sock path/to/config.toml
```

The server will start listening on the socket of your choice. You can then send SOAP requests to the configured endpoints.

## Thanks to
- [Teknoraver](https://github.com/teknoraver) for [this gist](https://gist.github.com/teknoraver/5ffacb8757330715bcbcc90e6d46ac74) about implementing http listen on Unix sockets
- [Foomo](https://github.com/foomo/) for their SOAP library, which I embedded in this project because I needed some downstream modifications

## Contributing

If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request on the GitHub repository.
