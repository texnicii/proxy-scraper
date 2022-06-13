package proxyscraper

import (
	"github.com/texnicii/proxy-scraper/parser/proxy"
	"reflect"
	"testing"
)

func TestCollectUniq(t *testing.T) {
	a := proxy.Proxy{
		Ipv4:      "0.0.0.0",
		Port:      0,
		ProxyType: "0",
	}
	b := proxy.Proxy{
		Ipv4:      "1.0.0.0",
		Port:      0,
		ProxyType: "0",
	}
	//is not uniq by ip address with "a"
	c := proxy.Proxy{
		Ipv4:      "0.0.0.0",
		Port:      0,
		ProxyType: "1",
	}

	type args struct {
		inputData []proxy.Proxy
		errorCh   <-chan error
		debug     bool
	}
	tests := []struct {
		name string
		args args
		want map[string]proxy.Proxy
	}{
		{
			name: "Two uniq elements write to channel",
			args: args{
				inputData: []proxy.Proxy{a, b},
				errorCh:   nil,
				debug:     false,
			},
			want: map[string]proxy.Proxy{
				"0.0.0.0": a,
				"1.0.0.0": b,
			},
		},
		{
			name: "Two uniq elements write to channel",
			args: args{
				inputData: []proxy.Proxy{a, c},
				errorCh:   nil,
				debug:     false,
			},
			want: map[string]proxy.Proxy{
				"0.0.0.0": a,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCh := make(chan []proxy.Proxy)
			go func() {
				inputCh <- tt.args.inputData
				close(inputCh)
			}()
			if got := CollectUniq(inputCh, nil, tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectUniq() = %v, want %v", got, tt.want)
			}
		})
	}
}
