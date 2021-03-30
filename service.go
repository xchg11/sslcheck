package main

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"

	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func checkSslWebhost() (map[string]int64, error) {
	var err error
	var conn *tls.Conn
	days_expiry := make(map[string]int64)
	for _, domain := range config.CheckSsl.Domains {
		now := time.Now().Unix()
		conn, err = tls.Dial("tcp", domain+":443", nil)
		if err != nil {
			return nil, err
		}
		err = conn.VerifyHostname(domain)
		if err != nil {
			return nil, err
		}
		expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
		days_expiry[domain] = ((expiry.Unix() - now) / 60 / 60 / 24)
		// log.Infof("Issuer: %s\nExpiry: %v\n", conn.ConnectionState().PeerCertificates[0].Issuer, expiry.Unix())
	}
	return days_expiry, err
}
func sendTgMsg(userid []string) {
	bot, err := tgbotapi.NewBotAPI(config.CheckSsl.Notify.Telegram.Apitoken)
	if err != nil {
		log.Error(err)
	}
	bot.Debug = config.CheckSsl.Debug
	if config.CheckSsl.Notify.Telegram.Debug {
		log.Infof("Authorized on account %s", bot.Self.UserName)
	}
	for _, userid := range config.CheckSsl.Notify.Telegram.Touser {
		i64, err := strconv.ParseInt(userid, 10, 64)
		if err != nil {
			log.Error("error parseint userid telegram ", err)
			return
		}
		msg := tgbotapi.NewMessage(i64, "privert!!!")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error(err)
		}
	}
}
