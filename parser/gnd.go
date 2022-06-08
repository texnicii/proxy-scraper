package parser

import (
	"encoding/json"
	"errors"
	"github.com/texnicii/proxy-scraper/parser/proxy"
	"strconv"
)

type GndParser struct {
	AbstractParser
}

func (parser *GndParser) Parse(content []byte) ([]proxy.Proxy, error) {
	proxyList := make([]proxy.Proxy, 0)
	tmpProxy := struct {
		Data []struct {
			Host      string   `json:"ip"`
			Port      string   `json:"port"`
			ProxyType []string `json:"protocols"`
		}
	}{}
	err := json.Unmarshal(content, &tmpProxy)
	for _, p := range tmpProxy.Data {
		if err != nil || !proxy.ValidateType(p.ProxyType[0]) {
			continue
		}
		port, _ := strconv.Atoi(p.Port)
		proxyList = append(proxyList, proxy.Proxy{
			Ipv4:      p.Host,
			Port:      port,
			ProxyType: p.ProxyType[0],
		})
	}
	if len(proxyList) == 0 {
		return nil, errors.New("json proxy list parse error")
	}
	return proxyList, nil
}
