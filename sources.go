package proxyscraper

import (
	"encoding/json"
	"os"
)

const DefaultFilename string = "proxy_sources.json"
const DefaultParserName string = "UniversalProxiesParser"

type ProxySource struct {
	Url        string `json:"url"`
	ProxyType  string `json:"type"`
	ParserName string `json:"parser"`
}

type ProxySourceList []ProxySource

func NewProxySourcesFromJson(jsonContent []byte) ProxySourceList {
	sources := ProxySourceList{}
	jsonErr := json.Unmarshal(jsonContent, &sources)
	if jsonErr != nil {
		panic("Json read error")
	}
	for k := range sources {
		if sources[k].ParserName == "" {
			sources[k].ParserName = DefaultParserName
		}
	}
	return sources
}

func NewProxySourcesFromFile() ProxySourceList {
	workingDir, _ := os.Getwd()
	content, err := os.ReadFile(workingDir + string(os.PathSeparator) + DefaultFilename)
	if err != nil {
		panic("Open source file error")
	}
	return NewProxySourcesFromJson(content)
}
