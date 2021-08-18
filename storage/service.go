package storage

import (
	// Standard Library Imports
	"strings"

	// Internal Imports
	"github.com/matthewhartstonge/sassy/api"
)

type signedPermissionSpec struct {
	Index      int
	APIVersion api.Version
}

func signedPermissionMap() map[string]signedPermissionSpec {
	return map[string]signedPermissionSpec{
		"r": {
			Index:      0,
			APIVersion: api.VAll,
		},
		"a": {
			Index:      1,
			APIVersion: api.VAll,
		},
		"c": {
			Index:      2,
			APIVersion: api.VAll,
		},
		"w": {
			Index:      3,
			APIVersion: api.VAll,
		},
		"d": {
			Index:      4,
			APIVersion: api.VAll,
		},
		"x": {
			Index:      5,
			APIVersion: api.V2019_12_12,
		},
		"y": {
			Index:      6,
			APIVersion: api.V2020_02_10,
		},
		"l": {
			Index:      7,
			APIVersion: api.VAll,
		},
		"t": {
			Index:      8,
			APIVersion: api.V2019_12_12,
		},
		"m": {
			Index:      9,
			APIVersion: api.V2020_02_10,
		},
		"e": {
			Index:      10,
			APIVersion: api.V2020_02_10,
		},
		"o": {
			Index:      11,
			APIVersion: api.V2020_02_10,
		},
		"p": {
			Index:      12,
			APIVersion: api.V2020_02_10,
		},
	}
}

func NewSignedPermissions() *SignedPermissions {
	return &SignedPermissions{
		permissions: [13]string{},
	}
}

type SignedPermissions struct {
	// permissions must be in the following order to comply to azure specifications: "racwdx(y?)ltmeop"
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#specifying-permissions
	permissions [13]string
}

// FromString parses a string and pushes in the signed permissions.
func (p *SignedPermissions) FromString(input string) *SignedPermissions {
	if p == nil {
		*p = *NewSignedPermissions()
	}

	// Grab each character of user input
	inputPermissions := strings.Split(input, "")

	// Add any valid user input permissions into the required order.
	spMap := signedPermissionMap()
	for _, permission := range inputPermissions {
		conformedPermission := strings.ToLower(permission)
		if spec, ok := spMap[conformedPermission]; ok {
			p.permissions[spec.Index] = conformedPermission
		}
	}

	return p
}

// GetKey returns the key of the param value to be added to the query string.
func (p SignedPermissions) GetKey() string {
	return "sp"
}

// GetValue returns the permissions string in the required order.
func (p SignedPermissions) GetValue(version api.Version) string {
	spMap := signedPermissionMap()
	out := [len(p.permissions)]string{}
	for _, value := range p.permissions {
		if spec, ok := spMap[value]; ok {
			if spec.APIVersion == api.VAll || spec.APIVersion <= version {
				out[spec.Index] = value
			}
		}
	}

	return strings.Replace(
		strings.Join(out[:], ""),
		" ", "", -1,
	)
}
