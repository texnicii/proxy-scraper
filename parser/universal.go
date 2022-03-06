package parser

import (
	"errors"
	"local/proxyscraper/parser/proxy"
	"regexp"
	"strconv"
)

// UniversalProxiesParser parse IPv4 and port from proxy lists via regexp
type UniversalProxiesParser struct {
	AbstractParser
}

func (parser *UniversalProxiesParser) Parse(content []byte) ([]proxy.Proxy, error) {
	proxyList := make([]proxy.Proxy, 0)
	exp := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d+)`)
	matched := exp.FindAllStringSubmatch(string(content), -1)
	if len(matched) == 0 {
		return nil, errors.New("parse error")
	}
	for _, submatch := range matched {
		port, convertErr := strconv.Atoi(submatch[2])
		if convertErr != nil {
			continue
		}
		proxyList = append(proxyList, proxy.NewProxy(submatch[1], port, parser.ProxyType))
	}
	return proxyList, nil
}
