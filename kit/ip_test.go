package kit

import (
	"net/netip"
	"testing"
)

func TestDemo(t *testing.T) {
	addr, err := netip.ParseAddr("192.168.1.1")
	t.Logf("addr: %v, err: %v", addr, err)
}

func TestIsContainIPWithScope(t *testing.T) {
	type args struct {
		ip      string
		ipScope string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should contain",
			args: args{ip: "192.168.1.2", ipScope: "192.168.1.1-192.168.3.1"},
			want: true,
		},
		{
			name: "should not contain",
			args: args{ip: "192.168.4.1", ipScope: "192.168.1.1-192.168.3.1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContainIPWithRange(tt.args.ip, tt.args.ipScope); got != tt.want {
				t.Errorf("IsContainIPWithRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
