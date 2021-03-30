package main

import (
	log "github.com/Sirupsen/logrus"
)

var (
	config     *data
	config_err error
)

func init() {
	config_err := readConfig("config.yaml")
	if config_err != nil {
		log.Fatal(config_err)
	}

}
func main() {
	//
	days, err := checkSslWebhost()
	if err != nil {
		log.Error(err)
	}
	for domain, day_expiry := range days {
		log.Info("aa: ", domain, " ", day_expiry)

	}
	//sendTgMsg([]string{"655412449"})
}
