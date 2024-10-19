package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/soheilrt/checkrr/checkrr"
	"github.com/soheilrt/checkrr/client"
	"github.com/soheilrt/checkrr/config"
)

func main() {
	config, err := config.LoadConfig(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	log.SetLevel(log.InfoLevel)
	l, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.WithError(err).Errorf("Error parsing log level")
	} else {
		log.SetLevel(l)
	}

	checkrrs := []*checkrr.CheckRR{}
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
