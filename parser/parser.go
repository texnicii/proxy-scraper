package parser

import (
	. "github.com/texnicii/proxy-scraper/parser/proxy"
)

type Parser interface {
	Parse(content []byte) ([]Proxy, error)
	SetType(proxyType string)
}

type AbstractParser struct {
	ProxyType string
}

func InitParser(name string, proxyType string) Parser {
	var parser Parser
	switch name {
	case "JsonParser":
		parser = new(JsonProxiesParser)
	case "GndParser":
		parser = new(GndParser)
	default:
		parser = new(UniversalProxiesParser)
	}
	parser.SetType(proxyType)
	return parser
}

func (p *AbstractParser) SetType(proxyType string) {
	p.ProxyType = proxyType
}
