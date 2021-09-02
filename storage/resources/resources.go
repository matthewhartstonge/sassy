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

package resources

import (
	// Standard Library Imports
	"net/url"
	"strings"
)

type Resource string

const (
	Container Resource = "c"
	Directory Resource = "d"
	Blob      Resource = "b"
)

const paramKey = "sr"

type Resources []Resource

func (r Resources) ToString() string {
	sr := ""
	for _, resource := range r {
		sr += string(resource)
	}

	return sr
}

func (r Resources) SetParam(params *url.Values) {
	if len(r) > 0 {
		params.Add(paramKey, r.ToString())
	}
}

func (r Resources) GetParam() (resources string) {
	if len(r) > 0 {
		values := &url.Values{}
		r.SetParam(values)

		resources = values.Encode()
	}

	return
}

func (r Resources) GetURLDecodedParam() (resources string) {
	if len(r) > 0 {
		resources, _ = url.QueryUnescape(r.GetParam())
	}

	return
}

func Parse(resources string) Resources {
	vMap := map[Resource]struct{}{
		Container: {},
		Directory: {},
		Blob:      {},
	}

	var sr Resources
	splitResources := strings.Split(resources, "")
	for _, service := range splitResources {
		check := Resource(service)
		if _, ok := vMap[check]; ok {
			sr = append(sr, check)
		}
	}

	return sr
}
