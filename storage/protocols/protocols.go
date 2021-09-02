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

package protocols

import (
	// Standard Library Imports
	"net/url"
	"strings"
)

type Protocol string

const (
	HTTP  Protocol = "http"
	HTTPS Protocol = "https"
)

const (
	numProtocols = 2
	paramKey     = "spr"
)

type Protocols struct {
	// hasValues tracks whether protocols have been added.
	hasValues bool
	// protocols must be in the order https,http
	protocols [numProtocols]Protocol
}

func defaultProtocols() Protocols {
	return Protocols{
		hasValues: true,
		protocols: [numProtocols]Protocol{HTTPS, HTTP},
	}
}

func (r Protocols) ToString() string {
	var out []string
	for _, protocol := range r.protocols {
		if protocol != "" {
			out = append(out, string(protocol))
		}
	}

	return strings.Join(out, ",")
}

func (r Protocols) SetParam(params *url.Values) {
	if r.hasValues {
		params.Add(paramKey, r.ToString())
	}
}

func (r Protocols) GetParam() (protocols string) {
	if r.hasValues {
		values := &url.Values{}
		r.SetParam(values)

		protocols = values.Encode()
	}

	return
}

func (r Protocols) GetURLDecodedParam() (protocols string) {
	if r.hasValues {
		protocols, _ = url.QueryUnescape(r.GetParam())
	}

	return
}

// Parse returns a valid set of protocols. Possible values are both HTTPS and
// HTTP (https,http) or HTTPS only (https).
// The default value is https,http.
func Parse(protocols string) (spr Protocols) {
	spr = Protocols{
		hasValues: false,
		protocols: [numProtocols]Protocol{},
	}

	splitProtocols := strings.Split(protocols, ",")
	if len(splitProtocols) == 1 && string(HTTP) == strings.ToLower(splitProtocols[0]) {
		// Note that HTTP only is not a permitted value.
		return defaultProtocols()
	}

	sprMap := protocolMap()
	for _, protocol := range splitProtocols {
		check := Protocol(strings.ToLower(protocol))
		if protocolIndex, ok := sprMap[check]; ok {
			spr.protocols[protocolIndex] = check
			if !spr.hasValues {
				spr.hasValues = true
			}
		}
	}

	return spr
}

func protocolMap() map[Protocol]int {
	return map[Protocol]int{
		HTTPS: 0,
		HTTP:  1,
	}
}
