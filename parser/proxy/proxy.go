package proxy

const (
	HTTP   = "http"
	HTTPS  = "https"
	SOCKS4 = "socks4"
	SOCKS5 = "socks5"
)

type Proxy struct {
	Ipv4      string
	Port      int
	ProxyType string
}

func NewProxy(ipv4 string, port int, proxyType string) Proxy {
	return Proxy{Ipv4: ipv4, Port: port, ProxyType: proxyType}
}

func ValidateType(proxyType string) bool {
	switch proxyType {
	case HTTP, HTTPS, SOCKS4, SOCKS5:
		return true
	default:
		return false
	}
}
