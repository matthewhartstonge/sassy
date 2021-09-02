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
	Latest = V2020_10_02

	V2020_10_02 SignedVersion = "2020-10-02"
	V2020_08_04 SignedVersion = "2020-08-04"
	V2020_02_10 SignedVersion = "2020-02-10"
	V2019_12_12 SignedVersion = "2019-12-12"
	V2015_04_05 SignedVersion = "2015-04-05"

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
		V2020_10_02: {},
		V2020_08_04: {},
		V2020_02_10: {},
		V2019_12_12: {},
		VAll:        {},
	}

	check := SignedVersion(version)
	if _, ok = vMap[check]; ok {
		return check, ok
	}

	// Default to "latest" for the latest level of functionality.
	return Latest, false
}
