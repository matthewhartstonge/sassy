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

// Package versions is required. It provides enums to specify the signed
// storage service version to use to authorize requests made with this account
// SAS.
//
// Must be set to version 2015-04-05 or later.
//
// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-account-sas#specifying-account-sas-parameters
package versions

import (
	// Standard Library Imports
	"net/url"
)

// SignedVersion specifies the signed storage service version to use to
// authorize requests made with this account SAS.
type SignedVersion string

const (
	Latest = V20201002

	V20201002 SignedVersion = "2020-10-02"
	V20200804 SignedVersion = "2020-08-04"
	V20200210 SignedVersion = "2020-02-10"
	V20191212 SignedVersion = "2019-12-12"
	V20150405 SignedVersion = "2015-04-05"

	// VAll is just a placeholder to delineate where a given function/property
	// is available in all API versions.
	VAll SignedVersion = "*"
)

const paramKey = "sv"

func (v SignedVersion) ToString() string {
	return string(v)
}

func (v SignedVersion) SetParam(params *url.Values) {
	if v != "" {
		params.Add(paramKey, v.ToString())
	}
}

func (v SignedVersion) GetParam() (signedVersion string) {
	if v != "" {
		values := &url.Values{}
		v.SetParam(values)

		signedVersion = values.Encode()
	}

	return
}

func (v SignedVersion) GetURLDecodedParam() (signedVersion string) {
	if v != "" {
		signedVersion, _ = url.QueryUnescape(v.GetParam())
	}

	return
}

// Parse returns an API Version from a given string. Defaults to latest.
func Parse(version string) (v SignedVersion, ok bool) {
	vMap := map[SignedVersion]struct{}{
		V20201002: {},
		V20200804: {},
		V20200210: {},
		V20191212: {},
		VAll:      {},
	}

	check := SignedVersion(version)
	if _, ok = vMap[check]; ok {
		return check, ok
	}

	// Default to "latest" for the latest level of functionality.
	return Latest, false
}
