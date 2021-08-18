package api

type SignedVersion string

const (
	Latest = V2020_10_02

	V2020_10_02 SignedVersion = "2020-10-02"
	V2020_08_04 SignedVersion = "2020-08-04"
	V2020_02_10 SignedVersion = "2020-02-10"
	V2019_12_12 SignedVersion = "2019-12-12"

	// VAll is just a placeholder to delineate where a given function/property
	// is available in all API versions.
	VAll SignedVersion = "*"
)

// ParseVersion returns an API Version.
func ParseVersion(version string) (v SignedVersion, ok bool) {
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
