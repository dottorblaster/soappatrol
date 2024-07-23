package main

import (
	"os"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"

	"github.com/dottorblaster/soappatrol/internal/soappatrol"
)

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

	var config soappatrol.Config
	if _, err := toml.Decode(string(configString), &config); err != nil {
		logger.Errorw("Unable to parse config")
	}

	logger.Infow("Unix SOAP server")

	os.Remove(os.Args[1])

	soappatrolServer := soappatrol.New(config, logger)

	err = soappatrolServer.ListenAndServe(os.Args[1])
	if err != nil {
		panic(err)
	}
}
