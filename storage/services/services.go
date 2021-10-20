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

package services

import (
	// Standard Library Imports
	"net/url"
	"strings"
)

type SignedService string

const (
	Blob  SignedService = "b"
	Queue SignedService = "q"
	Table SignedService = "t"
	File  SignedService = "f"
)

const paramKey = "ss"

type SignedServices []SignedService

func (s SignedServices) ToString() string {
	ss := ""
	for _, service := range s {
		ss += string(service)
	}

	return ss
}

func (s *SignedServices) SetParam(params *url.Values) {
	if s != nil && len(*s) > 0 {
		params.Add(paramKey, s.ToString())
	}
}

func (s SignedServices) GetParam() (services string) {
	if len(s) > 0 {
		values := &url.Values{}
		s.SetParam(values)

		services = values.Encode()
	}

	return
}

func (s SignedServices) GetURLDecodedParam() (services string) {
	if len(s) > 0 {
		services, _ = url.QueryUnescape(s.GetParam())
	}

	return
}

func Parse(services string) SignedServices {
	vMap := map[SignedService]struct{}{
		Blob:  {},
		Queue: {},
		Table: {},
		File:  {},
	}

	var ss SignedServices
	splitServices := strings.Split(services, "")
	for _, service := range splitServices {
		check := SignedService(service)
		if _, ok := vMap[check]; ok {
			ss = append(ss, check)
		}
	}

	return ss
}
