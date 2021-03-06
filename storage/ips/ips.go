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
	// Standard Library Imports
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
)

const (
	paramKey         = "sip"
	ipRangeSeparator = "-"
)

type SignedIP string

// String implements Stringer.
func (s SignedIP) String() string {
	return string(s)
}

func (s SignedIP) SetParam(params *url.Values) {
	if s != "" {
		params.Add(paramKey, s.String())
	}
}

func (s SignedIP) GetParam() (signedIPs string) {
	if s != "" {
		values := &url.Values{}
		s.SetParam(values)

		signedIPs = values.Encode()
	}

	return
}

func Parse(ips string) (sip SignedIP, ok bool) {
	// Attempt to separate a provided IP range.
	splitIPs := strings.Split(strings.TrimSpace(ips), ipRangeSeparator)
	switch len(splitIPs) {
	case 1:
		// Single IP
		ip := parseIPv4(splitIPs[0])
		if ip == nil {
			// Unable to parse IPv4 successfully, invalid IP.
			return
		}

		sip = SignedIP(ip.String())

	case 2:
		// IP Range
		rangeStart := parseIPv4(splitIPs[0])
		rangeEnd := parseIPv4(splitIPs[1])

		if rangeStart == nil || rangeEnd == nil {
			// One of the IPs have not been parsed successfully, invalid IP
			// range.
			return
		}

		if bytes.Compare(rangeStart, rangeEnd) == 1 {
			// range start is greater than range end, invalid IP range.
			return
		}

		sip = SignedIP(fmt.Sprintf(
			"%s-%s",
			rangeStart.String(),
			rangeEnd.String(),
		))

	default:
		// Multiple range separators provided, invalid IP range.
		return
	}

	return sip, true
}

// parseIPv4 checks parses an IPv4 address, returns nil if invalid.
func parseIPv4(ip string) net.IP {
	if parsedIP := net.ParseIP(strings.TrimSpace(ip)); parsedIP != nil {
		if ipv4 := parsedIP.To4(); ipv4 != nil {
			return ipv4
		}
	}

	return nil
}
