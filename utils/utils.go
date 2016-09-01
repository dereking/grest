package utils

import (
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("x-forwarded-for")
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || ip == "unknown" {
		//golang
		ipports := strings.Split(r.RemoteAddr, ":")
		if len(ipports) > 1 {
			ip = ipports[0]
		}
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("http_client_ip")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("Remote_addr") //nginx
	}

	// 如果是多级代理，那么取第一个ip为客户ip
	ips := strings.Split(ip, ",")
	if len(ips) > 1 {
		ip = ips[0]
	}

	return ip

}
