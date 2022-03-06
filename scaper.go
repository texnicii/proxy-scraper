package proxyscraper

import (
	"errors"
	"fmt"
	"io"
	"local/proxyscraper/parser"
	"local/proxyscraper/parser/proxy"
	"net/http"
	"sync"
)

const ChunkSize = 8

type ContentContainer struct {
	content []byte
	source  ProxySource
}

func Get(args ...bool) map[string]proxy.Proxy {
	var debug bool
	for _, debug = range args {
		break
	}
	s := NewProxySourcesFromFile()
	contentCh := make(chan ContentContainer)
	parsedProxyCh := make(chan []proxy.Proxy)
	errCh := make(chan error)
	go StartDownload(s, contentCh, errCh)
	go Parse(contentCh, parsedProxyCh, errCh, debug)
	// sync channels
	return CollectUniq(parsedProxyCh, errCh, debug)
}

// StartDownload split sources by chunks and run downloader
func StartDownload(sources ProxySourceList, contentCh chan<- ContentContainer, errCh chan<- error) {
	defer close(contentCh)
	chunkNum := 0
	sourcesLen := len(sources)
	if sourcesLen <= ChunkSize {
		AsyncDownload(sources, contentCh, errCh)
	} else {
		for i := range sources {
			if i%ChunkSize == 0 || i == sourcesLen-1 {
				chunkStartIndex := chunkNum * ChunkSize
				chunkSources := sources[chunkStartIndex : chunkStartIndex+ChunkSize]
				AsyncDownload(chunkSources, contentCh, errCh)
				chunkNum++
			}
		}
	}
}

// AsyncDownload async download from URL and send content to output channel
func AsyncDownload(sources ProxySourceList, contentCh chan<- ContentContainer, errCh chan<- error) {
	wg := new(sync.WaitGroup)
	wg.Add(len(sources))
	for _, proxySource := range sources {
		go func(source ProxySource, group *sync.WaitGroup, outputContent chan<- ContentContainer, outputErr chan<- error) {
			defer group.Done()
			if source.Url == "" {
				return
			}
			response, err := http.Get(source.Url)
			if err != nil {
				errCh <- err
				return
			}
			content, _ := io.ReadAll(response.Body)
			if len(content) == 0 {
				errCh <- errors.New("empty content (" + source.Url + ")")
				return
			}
			outputContent <- ContentContainer{content, source}
		}(proxySource, wg, contentCh, errCh)
	}
	wg.Wait()
}

func Parse(contentCh <-chan ContentContainer, output chan<- []proxy.Proxy, errCh chan<- error, verbose bool) {
	defer close(output)
	defer close(errCh)
	for {
		if contentContainer, ok := <-contentCh; ok {
			proxyParser := parser.InitParser(contentContainer.source.ParserName, contentContainer.source.ProxyType)
			proxyList, parseErr := proxyParser.Parse(contentContainer.content)
			if parseErr != nil {
				errCh <- parseErr
				continue
			}
			if verbose {
				fmt.Println(contentContainer.source.Url, len(proxyList))
			}
			output <- proxyList
		} else {
			break
		}
	}
}

// CollectUniq makes unique proxy by IP
func CollectUniq(input <-chan []proxy.Proxy, errorCh <-chan error, debug bool) map[string]proxy.Proxy {
	proxies := make(map[string]proxy.Proxy)
out:
	for {
		select {
		case proxyListChunk, ok := <-input:
			if !ok {
				break out
			}
			for _, proxyItem := range proxyListChunk {
				if _, ok := proxies[proxyItem.Ipv4]; !ok {
					proxies[proxyItem.Ipv4] = proxyItem
				}
			}
		case err := <-errorCh:
			if debug && err != nil {
				fmt.Println(err)
			}
		}
	}
	return proxies
}
