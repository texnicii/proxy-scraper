package main

import (
	"fmt"
	"local/proxyscraper"
	"os"
	"time"
)

func main() {
	var debug bool
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		debug = true
	}
	proxies := proxyscraper.Get(debug)
	filename := "proxy_" + time.Now().Format("2006-01-02-150405") + ".csv"
	outputFile, openErr := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if openErr != nil {
		return
	}
	for _, proxy := range proxies {
		_, writeErr := outputFile.WriteString(proxy.String())
		if writeErr != nil {
			return
		}
	}
	fmt.Printf("Found proxy: %d\nStored in: %s\n", len(proxies), filename)
}
