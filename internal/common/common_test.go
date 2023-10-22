package common

import "testing"

func TestExpandIPv6Zero(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "::",
			args: args{
				ip: "::",
			},
			want: "0:0:0:0:0:0:0:0",
		},
		{
			name: "::1",
			args: args{
				ip: "::1",
			},
			want: "0:0:0:0:0:0:0:1",
		},
		{
			name: "1::",
			args: args{
				ip: "1::",
			},
			want: "1:0:0:0:0:0:0:0",
		},
		{
			name: "1::1",
			args: args{
				ip: "1::1",
			},
			want: "1:0:0:0:0:0:0:1",
		},
		{
			name: "2001:db8:85a3::8a2e:370:7334",
			args: args{
				ip: "2001:db8:85a3::8a2e:370:7334",
			},
			want: "2001:db8:85a3:0:0:8a2e:370:7334",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandIPv6Zero(tt.args.ip); got != tt.want {
				t.Errorf("ExpandIPv6Zero() = %v, want %v", got, tt.want)
			}
		})
	}
}
