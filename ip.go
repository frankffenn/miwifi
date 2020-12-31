package main

import (
	"regexp"
)

const (
	targetAPI = "http://api.89ip.cn/tqdl.html?api=1&num=100&port=&address=&isp="
	patten    = `((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`
)

func fetchIP() ([]string, error) {
	resp, err := HTTPGet(targetAPI)
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(patten)
	out := reg.FindAllString(string(resp), -1)

	return out, nil
}
