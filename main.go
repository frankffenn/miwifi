package main

import (
	"log"
	"time"
)

func main() {
	m := NewMiWifi()
	resp, _ := m.Call()
	if resp.Code < 0 {
		log.Println("auth failed:", resp.Message)
		return
	}

	log.Println("connect wifi success, check your internet!")

	ips, err := fetchIP()
	if err != nil {
		log.Println("fetch ip failed:%v", err)
		return
	}

	for _, ip := range ips {

		m.SetProxy(ip)

		for i := 0; i < 4; i++ {
			time.Sleep(100 * time.Second)

			resp, _ := m.Call()
			if resp.Code < 0 {
				log.Println("auth failed:", resp.Message)
				return
			}

			log.Printf("keep alive ..., ip:%s, retry times:%d\n", ip, i+1)
		}
	}
}
