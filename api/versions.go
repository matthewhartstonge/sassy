package api

type Version string

const (
	V2020_10_02 Version = "2020-10-02"
	V2020_08_04 Version = "2020-08-04"
	V2020_02_10 Version = "2020-02-10"
	V2019_12_12 Version = "2019-12-12"

	// VAll is just a placeholder to delineate where a given function/property
	// is available in all API versions.
	VAll Version = "*"
)

// ParseVersion returns an API Version.
func ParseVersion(version string) (v Version, ok bool) {
	vMap := map[Version]struct{}{
		V2020_10_02: {},
		V2020_08_04: {},
		V2020_02_10: {},
		V2019_12_12: {},
		VAll:        {},
	}

	check := Version(version)
	if _, ok = vMap[check]; ok {
		return check, ok
	}

	// Default to "vall" for the oldest level of compatibility.
	return VAll, false
}
