package proxy

import "testing"

func TestProxy_GetAddress(t *testing.T) {
	type fields struct {
		Ipv4      string
		Port      int
		ProxyType string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Make address as string ip:port",
			fields: fields{
				Ipv4:      "0.0.0.0",
				Port:      1234,
				ProxyType: "https",
			},
			want: "0.0.0.0:1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Proxy{
				Ipv4:      tt.fields.Ipv4,
				Port:      tt.fields.Port,
				ProxyType: tt.fields.ProxyType,
			}
			if got := p.GetAddress(); got != tt.want {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
