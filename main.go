package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/soheilrt/checkrr/pkg/checkrr"
	"github.com/soheilrt/checkrr/pkg/client"
	"github.com/soheilrt/checkrr/pkg/config"
)

func main() {
	configFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return
	}
	defer configFile.Close()

	config, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	log.SetLevel(log.InfoLevel)
	l, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.WithError(err).Errorf("Error parsing log level")
	} else {
		log.SetLevel(l)
	}

	var checkrrs []*checkrr.CheckRR
	for _, c := range config.Clients {
		cl := client.NewClientRR(c.Host, c.APIKey, c.Options)
		check := checkrr.NewCheckRR(c.Name, cl, c.Conditions)
		checkrrs = append(checkrrs, check)
	}

	for {
		for _, c := range checkrrs {
			err := c.Check()
			if err != nil {
				log.WithError(err).Errorf("error checking client")
			}
		}
		time.Sleep(config.SleepTime)
	}
}
