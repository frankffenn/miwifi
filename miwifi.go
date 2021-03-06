package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	SnsInit      = "http://guest.miwifi.com:8999/cgi-bin/luci/api/misns/sns_init"
	ProtalConfig = "http://api.miwifi.com/guest_wifi/portal_config?callback=jQuery210032440425060714295_%d&did=%s&client_info=%s&_=%d"
	ApplyRent    = "http://api.miwifi.com/wifirent/api/ad_apply_rent?callback=jsonpCallback&router_id=%s&client_info=%s&_=%d"
)

type MiWifi struct {
	DeviceID   string
	ClientInfo string
	Proxy      string
}

func NewMiWifi() *MiWifi {
	return &MiWifi{}
}

type Sns struct {
	Code       int    `json:"code"`
	SSID       string `json:"ssid"`
	DeviceID   string `json:"deviceid"`
	ClientInfo string `json:"clientinfo"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (m *MiWifi) SetProxy(proxy string) {
	m.Proxy = proxy
}

func (m *MiWifi) init() (*Sns, error) {
	log.Println("[HTTP]", "GET", SnsInit)
	resp, err := m.doRequest(SnsInit)
	if err != nil {
		return nil, err
	}
	resp = resp[13 : len(resp)-2]
	log.Println(string(resp))

	var out Sns
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	m.DeviceID = out.DeviceID
	m.ClientInfo = out.ClientInfo
	return &out, nil
}

func (m *MiWifi) protalConfig(deviceId, clientInfo string) error {
	times := time.Now().Unix() - 2000
	log.Println("[HTTP]", "GET", fmt.Sprintf(ProtalConfig, times, deviceId, clientInfo, times+2))
	_, err := m.doRequest(fmt.Sprintf(ProtalConfig, times, deviceId, clientInfo, times+2))
	if err != nil {
		return err
	}
	return nil
}

func (m *MiWifi) apply(deviceId, clientInfo string) (*Response, error) {
	times := time.Now().Unix() - 2000
	log.Println("[HTTP]", "GET", fmt.Sprintf(ApplyRent, deviceId, clientInfo, times))
	resp, err := m.doRequest(fmt.Sprintf(ApplyRent, deviceId, clientInfo, times))
	if err != nil {
		return nil, err
	}
	resp = resp[14 : len(resp)-2]
	log.Println(string(resp))

	var out Response
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (m *MiWifi) doRequest(url string) ([]byte, error) {
	if m.Proxy == "" {
		return HTTPGet(url)
	}
	return HTTPGetWithProxy(url, m.Proxy)
}

func (m *MiWifi) Call() (*Response, error) {
	sns, err := m.init()
	if err != nil {
		log.Println("sns_init error", "err", err)
		return nil, err
	}

	if err := m.protalConfig(sns.DeviceID, sns.ClientInfo); err != nil {
		log.Println("protal config error", "err", err)
		return nil, err
	}

	resp, err := m.apply(sns.DeviceID, sns.ClientInfo)
	if err != nil {
		log.Println("apply rent error", "err", err)
		return nil, err
	}

	return resp, nil
}

func (m *MiWifi) KeepAlive() {
	req, _ := http.NewRequest(http.MethodGet, "http://www.baidu.com", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Minute,
				KeepAlive: 120 * time.Minute,
			}).DialContext,
			MaxConnsPerHost: 0,
		},
	}

	for {
		_, err := client.Do(req)
		if err != nil {
			log.Println("err", err)
		}

		log.Println("keep alive ...")

		select {
		case <-time.After(5 * time.Second):
		}
	}
}
