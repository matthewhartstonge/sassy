package signedip

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

func (p SignedIP) ToString() string {
	return string(p)
}

func (p SignedIP) SetParam(params *url.Values) {
	if p != "" {
		params.Add(paramKey, p.ToString())
	}
}

func (p SignedIP) GetParam() (signedIPs string) {
	if p != "" {
		values := &url.Values{}
		p.SetParam(values)

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
