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

type ResourceType string

const (
	Service   ResourceType = "s"
	Container ResourceType = "c"
	Object    ResourceType = "o"
)

const paramKey = "srt"

type ResourceTypes []ResourceType

func (r ResourceTypes) ToString() string {
	sr := ""
	for _, resource := range r {
		sr += string(resource)
	}

	return sr
}

func (r ResourceTypes) SetParam(params *url.Values) {
	if len(r) > 0 {
		params.Add(paramKey, r.ToString())
	}
}

func (r ResourceTypes) GetParam() (resourceTypes string) {
	if len(r) > 0 {
		values := &url.Values{}
		r.SetParam(values)

		resourceTypes = values.Encode()
	}

	return
}

func (r ResourceTypes) GetURLDecodedParam() (resourceTypes string) {
	if len(r) > 0 {
		resourceTypes, _ = url.QueryUnescape(r.GetParam())
	}

	return
}

func Parse(resourceTypes string) (srt ResourceTypes) {
	srtMap := resourceTypeMap()
	splitResourceTypes := strings.Split(resourceTypes, "")
	for _, resourceType := range splitResourceTypes {
		check := ResourceType(resourceType)
		if _, ok := srtMap[check]; ok {
			srt = append(srt, check)
		}
	}

	return srt
}

func resourceTypeMap() map[ResourceType]struct{} {
	return map[ResourceType]struct{}{
		Service:   {},
		Container: {},
		Object:    {},
	}
}
