package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/texnicii/proxy-scraper/parser/proxy"
)

type JsonProxiesParser struct {
	AbstractParser
}

func (parser *JsonProxiesParser) Parse(content []byte) ([]proxy.Proxy, error) {
	proxyList := make([]proxy.Proxy, 0)
	lines := bytes.Split(content, []byte("\n"))
	for _, jsonLine := range lines {
		tmpProxy := struct {
			Host      string `json:"host"`
			Port      int    `json:"port"`
			ProxyType string `json:"type"`
		}{}
		err := json.Unmarshal(jsonLine, &tmpProxy)
		if err != nil || !proxy.ValidateType(tmpProxy.ProxyType) {
			continue
		}
		proxyList = append(proxyList, proxy.Proxy{
			Ipv4:      tmpProxy.Host,
			Port:      tmpProxy.Port,
			ProxyType: tmpProxy.ProxyType,
		})
	}
	if len(proxyList) == 0 {
		return nil, errors.New("json proxy list parse error")
	}
	return proxyList, nil
}
