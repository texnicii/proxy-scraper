package proxy

import (
	"bytes"
	"encoding/csv"
	"strconv"
)

const (
	HTTP   = "http"
	HTTPS  = "https"
	SOCKS4 = "socks4"
	SOCKS5 = "socks5"
)

type Proxy struct {
	Ipv4      string
	Port      int
	ProxyType string
}

func NewProxy(ipv4 string, port int, proxyType string) Proxy {
	return Proxy{Ipv4: ipv4, Port: port, ProxyType: proxyType}
}

func ValidateType(proxyType string) bool {
	switch proxyType {
	case HTTP, HTTPS, SOCKS4, SOCKS5:
		return true
	default:
		return false
	}
}

func (p Proxy) String() string {
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	err := csvWriter.Write([]string{p.Ipv4, strconv.Itoa(p.Port), p.ProxyType})
	if err != nil {
		return ""
	}
	csvWriter.Flush()
	return buf.String()
}
