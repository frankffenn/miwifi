package main

import (
	"log"
)

func main() {

	m := NewMiWifi()
	sns, err := m.Init()
	if err != nil {
		log.Println("sns_init error", "err", err)
		return
	}

	if err := m.ProtalConfig(sns.DeviceID, sns.ClientInfo); err != nil {
		log.Println("protal config error", "err", err)
		return
	}

	resp, err := m.Apply(sns.DeviceID, sns.ClientInfo)
	if err != nil {
		log.Println("apply rent error", "err", err)
		return
	}

	if resp.Code < 0 {
		log.Println("auth failed:", resp.Message)
		return
	}

	log.Println("connect wifi success, check your internet!")

	m.KeepAlive()
}
