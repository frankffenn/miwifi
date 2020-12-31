package main

import (
	"log"
)

func main() {
	m := NewMiWifi()
	ips, err := fetchIP()
	if err != nil {
		log.Println("fetch ip failed:%v", err)
		return
	}

	for _, ip := range ips {

		m.SetProxy(ip)

		for i := 0; i < 4; i++ {
			resp, _ := m.Call()
			if resp.Code < 0 {
				log.Println("auth failed:", resp.Message)
				return
			}

			log.Println("connect wifi success, check your internet!")
		}
	}

	m.KeepAlive()
}
