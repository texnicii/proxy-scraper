package proxyscraper

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Pagination struct {
		Start   int    `json:"start"`
		End     int    `json:"end"`
		Pattern string `json:"pattern"`
	}
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
		//source has a pages
		if sources[k].Pagination.Pattern != "" {
			for i := sources[k].Pagination.Start; i <= sources[k].Pagination.End; i++ {
				copiedSource := sources[k]
				copiedSource.Url = copiedSource.Url + fmt.Sprintf(copiedSource.Pagination.Pattern, i)
				sources = append(sources, copiedSource)
			}
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
