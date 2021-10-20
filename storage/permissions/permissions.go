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

package permissions

import (
	// Standard Library Imports
	"net/url"
	"strings"

	// Internal Imports
	"github.com/matthewhartstonge/sassy/storage/versions"
)

const (
	numPermissions = 13
	paramKey       = "sp"
)

type SignedPermission string

// String implements Stringer.
func (s SignedPermission) String() string {
	return string(s)
}

const (
	Read            SignedPermission = "r"
	Add             SignedPermission = "a"
	Create          SignedPermission = "c"
	Write           SignedPermission = "w"
	Delete          SignedPermission = "d"
	DeleteVersion   SignedPermission = "x"
	PermanentDelete SignedPermission = "y"
	List            SignedPermission = "l"
	Tags            SignedPermission = "t"
	Move            SignedPermission = "m"
	Execute         SignedPermission = "e"
	Ownership       SignedPermission = "o"
	Permissions     SignedPermission = "p"
)

type SignedPermissions struct {
	// in order to return/parse the correct permissions, we need to know the
	// API version in use.
	versions.SignedVersion

	// hasValues tracks whether permissions have been added.
	hasValues bool
	// permissions must be in the following order to comply to azure
	// specifications: "racwdxyltmeop"
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#specifying-permissions
	permissions [numPermissions]SignedPermission
}

func (s SignedPermissions) String() string {
	var out []string
	spMap := signedPermissionMap()
	for _, permission := range s.permissions {
		if spec, ok := spMap[permission]; ok {
			if spec.APIVersion == versions.VAll || spec.APIVersion <= s.SignedVersion {
				out = append(out, permission.String())
			}
		}
	}

	return strings.Join(out, "")
}

func (s SignedPermissions) SetParam(params *url.Values) {
	if s.hasValues {
		params.Add(paramKey, s.String())
	}
}

func (s SignedPermissions) GetParam() (signedPermissions string) {
	if s.hasValues {
		values := &url.Values{}
		s.SetParam(values)

		signedPermissions = values.Encode()
	}

	return
}

func (s SignedPermissions) GetURLDecodedParam() (signedPermissions string) {
	if s.hasValues {
		signedPermissions, _ = url.QueryUnescape(s.GetParam())
	}

	return
}

func Parse(version versions.SignedVersion, permissions string) (sp SignedPermissions) {
	sp = SignedPermissions{
		hasValues:   false,
		permissions: [numPermissions]SignedPermission{},

		SignedVersion: version,
	}

	spMap := signedPermissionMap()
	splitPermissions := strings.Split(strings.ToLower(strings.TrimSpace(permissions)), "")
	for _, permission := range splitPermissions {
		signedPermission := SignedPermission(permission)
		if spec, ok := spMap[signedPermission]; ok {
			sp.permissions[spec.Index] = signedPermission
			if !sp.hasValues {
				sp.hasValues = true
			}
		}
	}

	return sp
}

type signedPermissionSpec struct {
	OpName        string
	OpDescription string
	Index         int
	APIVersion    versions.SignedVersion
}

func signedPermissionMap() map[SignedPermission]signedPermissionSpec {
	// nextIndex provides an index generating closure, so we don't have to
	// manually track indices in the map.
	i := -1
	nextIndex := func() int {
		i++
		return i
	}

	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#permissions-for-a-directory-container-or-blob
	return map[SignedPermission]signedPermissionSpec{
		Read: {
			OpName:        "Read",
			OpDescription: "Read the content, block list, properties, and metadata of any blob in the container or directory. Use a blob as the source of a copy operation.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		Add: {
			OpName:        "Add",
			OpDescription: "Add a block to an append blob.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		Create: {
			OpName:        "Create",
			OpDescription: "Write a new blob, snapshot a blob, or copy a blob to a new blob.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		Write: {
			OpName:        "Write",
			OpDescription: "Create or write content, properties, metadata, or block list. Snapshot or lease the blob. Resize the blob (page blob only). Use the blob as the destination of a copy operation.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		Delete: {
			OpName:        "Delete",
			OpDescription: "Delete a blob. For version 2017-07-29 and later, the Delete permission also allows breaking a lease on a blob. For more information, see the Lease Blob operation.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		DeleteVersion: {
			OpName:        "Delete version",
			OpDescription: "Delete a blob version.",
			Index:         nextIndex(),
			APIVersion:    versions.V20191212,
		},
		PermanentDelete: {
			OpName:        "Permanent delete",
			OpDescription: "Permanently delete a blob snapshot or version.",
			Index:         nextIndex(),
			APIVersion:    versions.V20200210,
		},
		List: {
			OpName:        "List",
			OpDescription: "List blobs non-recursively.",
			Index:         nextIndex(),
			APIVersion:    versions.VAll,
		},
		Tags: {
			OpName:        "Tags",
			OpDescription: "Read or write the tags on a blob.",
			Index:         nextIndex(),
			APIVersion:    versions.V20191212,
		},
		Move: {
			OpName:        "Move",
			OpDescription: "Move a blob or a directory and its contents to a new location. This operation can optionally be restricted to the owner of the child blob, directory, or parent directory if the `saoid` parameter is included on the SAS token and the sticky bit is set on the parent directory.",
			Index:         nextIndex(),
			APIVersion:    versions.V20200210,
		},
		Execute: {
			OpName:        "Execute",
			OpDescription: "Get the system properties and, if the hierarchical namespace is enabled for the storage account, get the POSIX ACL of a blob. If the hierarchical namespace is enabled and the caller is the owner of a blob, this permission grants the ability to set the owning group, POSIX permissions, and POSIX ACL of the blob. Does not permit the caller to read user-defined metadata.",
			Index:         nextIndex(),
			APIVersion:    versions.V20200210,
		},
		Ownership: {
			OpName:        "Ownership",
			OpDescription: "When the hierarchical namespace is enabled, this permission enables the caller to set the owner or the owning group, or to act as the owner when renaming or deleting a directory or blob within a directory that has the sticky bit set.",
			Index:         nextIndex(),
			APIVersion:    versions.V20200210,
		},
		Permissions: {
			OpName:        "Permissions",
			OpDescription: "When the hierarchical namespace is enabled, this permission allows the caller to set permissions and POSIX ACLs on directories and blobs.",
			Index:         nextIndex(),
			APIVersion:    versions.V20200210,
		},
	}
}
