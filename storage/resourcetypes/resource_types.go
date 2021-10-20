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

package resourcetypes

import (
	// Standard Library Imports
	"net/url"
	"strings"
)

type SignedResourceType string

// String implements Stringer.
func (s SignedResourceType) String() string {
	return string(s)
}

const (
	Service   SignedResourceType = "s"
	Container SignedResourceType = "c"
	Object    SignedResourceType = "o"
)

const paramKey = "srt"

type SignedResourceTypes []SignedResourceType

func (s SignedResourceTypes) String() string {
	sr := ""
	for _, resource := range s {
		sr += resource.String()
	}

	return sr
}

func (s SignedResourceTypes) SetParam(params *url.Values) {
	if len(s) > 0 {
		params.Add(paramKey, s.String())
	}
}

func (s SignedResourceTypes) GetParam() (resourceTypes string) {
	if len(s) > 0 {
		values := &url.Values{}
		s.SetParam(values)

		resourceTypes = values.Encode()
	}

	return
}

func (s SignedResourceTypes) GetURLDecodedParam() (resourceTypes string) {
	if len(s) > 0 {
		resourceTypes, _ = url.QueryUnescape(s.GetParam())
	}

	return
}

func Parse(resourceTypes string) (srt SignedResourceTypes) {
	srtMap := resourceTypeMap()
	splitResourceTypes := strings.Split(strings.ToLower(strings.TrimSpace(resourceTypes)), "")
	for _, resourceType := range splitResourceTypes {
		check := SignedResourceType(resourceType)
		if _, ok := srtMap[check]; ok {
			srt = append(srt, check)
		}
	}

	return srt
}

func resourceTypeMap() map[SignedResourceType]struct{} {
	return map[SignedResourceType]struct{}{
		Service:   {},
		Container: {},
		Object:    {},
	}
}
