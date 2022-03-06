package main

import (
	"fmt"
	"local/proxyscraper"
)

func main() {
	proxies := proxyscraper.Get(true)
	fmt.Printf("Found proxy: %d\n", len(proxies))
}
