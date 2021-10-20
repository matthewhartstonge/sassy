/*
 * Copyright © 2021 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ips

import (
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		ips string
	}
	tests := []struct {
		name    string
		args    args
		wantSip SignedIP
		wantOk  bool
	}{
		// Single IP testing
		{
			name: "Should not allow a partial IPv4 address",
			args: args{
				ips: "1.1",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow an invalid IPv4 address",
			args: args{
				ips: "1.1.1.256",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should allow a valid IPv4",
			args: args{
				ips: "1.1.1.1",
			},
			wantSip: "1.1.1.1",
			wantOk:  true,
		},

		// Disallow IPv6
		{
			name: "Should not allow an IPv6 address",
			args: args{
				ips: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow an IPv6 address range",
			args: args{
				ips: "2001:0db8:85a3:0000:0000:8a2e:0370:7334-2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			wantSip: "",
			wantOk:  false,
		},

		// Range Testing
		{
			name: "Should not allow a range containing a partial starting IPv4",
			args: args{
				ips: "1.1-1.1.1.1",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow a range containing an invalid starting IPv4",
			args: args{
				ips: "1.1.1.256-1.1.1.1",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow a range containing a partial ending IPv4",
			args: args{
				ips: "1.1.1.1-1.1",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow a range containing an invalid ending IPv4",
			args: args{
				ips: "1.1.1.1-1.1.1.256",
			},
			wantSip: "",
			wantOk:  false,
		},

		// Range testing for First Octet
		{
			name: "Should not allow an IPv4 address where first octet is lower on end",
			args: args{
				ips: "2.2.2.2-1.2.2.2",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should allow an IPv4 address where first octet matches on end",
			args: args{
				ips: "1.2.2.2-1.2.2.2",
			},
			wantSip: "1.2.2.2-1.2.2.2",
			wantOk:  true,
		},
		{
			name: "Should allow an IPv4 address where first octet is higher on end",
			args: args{
				ips: "2.2.2.2-3.2.2.2",
			},
			wantSip: "2.2.2.2-3.2.2.2",
			wantOk:  true,
		},

		// Range testing for Second Octet
		{
			name: "Should not allow an IPv4 address where second octet is lower on end",
			args: args{
				ips: "2.2.2.2-2.1.2.2",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should allow an IPv4 address where second octet matches on end",
			args: args{
				ips: "2.1.2.2-2.1.2.2",
			},
			wantSip: "2.1.2.2-2.1.2.2",
			wantOk:  true,
		},
		{
			name: "Should allow an IPv4 address where second octet is higher on end",
			args: args{
				ips: "2.2.2.2-2.3.2.2",
			},
			wantSip: "2.2.2.2-2.3.2.2",
			wantOk:  true,
		},

		// Range testing for Third Octet
		{
			name: "Should not allow an IPv4 address where third octet is lower on end",
			args: args{
				ips: "2.2.2.2-2.2.1.2",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should allow an IPv4 address where third octet matches on end",
			args: args{
				ips: "2.2.1.2-2.2.1.2",
			},
			wantSip: "2.2.1.2-2.2.1.2",
			wantOk:  true,
		},
		{
			name: "Should allow an IPv4 address where third octet is higher on end",
			args: args{
				ips: "2.2.2.2-2.2.3.2",
			},
			wantSip: "2.2.2.2-2.2.3.2",
			wantOk:  true,
		},

		// Range testing for Fourth Octet
		{
			name: "Should not allow an IPv4 address where fourth octet is lower on end",
			args: args{
				ips: "2.2.2.2-2.2.2.1",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should allow an IPv4 address where fourth octet matches on end",
			args: args{
				ips: "2.2.2.1-2.2.2.1",
			},
			wantSip: "2.2.2.1-2.2.2.1",
			wantOk:  true,
		},
		{
			name: "Should allow an IPv4 address where fourth octet is higher on end",
			args: args{
				ips: "2.2.2.2-2.2.2.3",
			},
			wantSip: "2.2.2.2-2.2.2.3",
			wantOk:  true,
		},

		{
			name: "Should not allow double separation between ranges",
			args: args{
				ips: "2.2.2.2--2.2.2.3",
			},
			wantSip: "",
			wantOk:  false,
		},
		{
			name: "Should not allow multiple ranges",
			args: args{
				ips: "2.2.2.2-2.2.2.3-2.2.2.4",
			},
			wantSip: "",
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSip, gotOk := Parse(tt.args.ips)
			if gotSip != tt.wantSip {
				t.Errorf("Parse() sip\ngot:  = %v\nwant: %v\n", gotSip, tt.wantSip)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Parse() ok\ngot:  = %v\nwant: %v\n", gotOk, tt.wantOk)
			}
		})
	}
}
