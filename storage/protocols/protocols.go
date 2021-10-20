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

type SignedProtocol string

const (
	HTTP  SignedProtocol = "http"
	HTTPS SignedProtocol = "https"
)

const (
	numProtocols = 2
	paramKey     = "spr"
)

type SignedProtocols struct {
	// hasValues tracks whether protocols have been added.
	hasValues bool
	// protocols must be in the order https,http
	protocols [numProtocols]SignedProtocol
}

func defaultProtocols() SignedProtocols {
	return SignedProtocols{
		hasValues: true,
		protocols: [numProtocols]SignedProtocol{HTTPS, HTTP},
	}
}

func (r SignedProtocols) ToString() string {
	var out []string
	for _, protocol := range r.protocols {
		if protocol != "" {
			out = append(out, string(protocol))
		}
	}

	return strings.Join(out, ",")
}

func (r SignedProtocols) SetParam(params *url.Values) {
	if r.hasValues {
		params.Add(paramKey, r.ToString())
	}
}

func (r SignedProtocols) GetParam() (protocols string) {
	if r.hasValues {
		values := &url.Values{}
		r.SetParam(values)

		protocols = values.Encode()
	}

	return
}

func (r SignedProtocols) GetURLDecodedParam() (protocols string) {
	if r.hasValues {
		protocols, _ = url.QueryUnescape(r.GetParam())
	}

	return
}

// Parse returns a valid set of protocols. Possible values are both HTTPS and
// HTTP (https,http) or HTTPS only (https).
// The default value is https,http.
func Parse(protocols string) (spr SignedProtocols) {
	spr = SignedProtocols{
		hasValues: false,
		protocols: [numProtocols]SignedProtocol{},
	}

	splitProtocols := strings.Split(protocols, ",")
	if len(splitProtocols) == 1 && string(HTTP) == strings.ToLower(splitProtocols[0]) {
		// Note that HTTP only is not a permitted value.
		return defaultProtocols()
	}

	sprMap := protocolMap()
	for _, protocol := range splitProtocols {
		check := SignedProtocol(strings.ToLower(protocol))
		if protocolIndex, ok := sprMap[check]; ok {
			spr.protocols[protocolIndex] = check
			if !spr.hasValues {
				spr.hasValues = true
			}
		}
	}

	return spr
}

func protocolMap() map[SignedProtocol]int {
	return map[SignedProtocol]int{
		HTTPS: 0,
		HTTP:  1,
	}
}
