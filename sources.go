package proxyscraper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const DefaultFilename string = "proxy_sources.json"
const DefaultRemoteSourcesURL string = "https://raw.githubusercontent.com/texnicii/proxy-scraper/master/proxy_sources.json"
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
		panic(jsonErr)
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
	localSourcesFilename := workingDir + string(os.PathSeparator) + DefaultFilename
	var content []byte
	var err error
	if _, err = os.Stat(localSourcesFilename); err == nil {
		content, err = os.ReadFile(localSourcesFilename)
	} else {
		var response *http.Response
		response, err = http.Get(DefaultRemoteSourcesURL)
		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(response.Body)
		if err == nil && response.StatusCode == 200 {
			content, err = io.ReadAll(response.Body)
		} else {
			err = errors.New("download sources file error")
		}
	}
	if err != nil {
		panic(err)
	}
	return NewProxySourcesFromJson(content)
}
