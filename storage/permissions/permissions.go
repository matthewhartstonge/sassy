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

type SignedPermissions struct {
	// in order to return/parse the correct permissions, we need to know the
	// API version in use.
	versions.SignedVersion

	// hasValues tracks whether permissions have been added.
	hasValues bool
	// permissions must be in the following order to comply to azure
	// specifications: "racwdxyltmeop"
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#specifying-permissions
	permissions [numPermissions]string
}

func (p SignedPermissions) ToString() string {
	var out []string
	spMap := signedPermissionMap()
	for _, value := range p.permissions {
		if spec, ok := spMap[value]; ok {
			if spec.APIVersion == versions.VAll || spec.APIVersion <= p.SignedVersion {
				out = append(out, value)
			}
		}
	}

	return strings.Join(out, "")
}

func (p SignedPermissions) SetParam(params *url.Values) {
	if p.hasValues {
		params.Add(paramKey, p.ToString())
	}
}

func (p SignedPermissions) GetParam() (signedPermissions string) {
	if p.hasValues {
		values := &url.Values{}
		p.SetParam(values)

		signedPermissions = values.Encode()
	}

	return
}

func (p SignedPermissions) GetURLDecodedParam() (signedPermissions string) {
	if p.hasValues {
		signedPermissions, _ = url.QueryUnescape(p.GetParam())
	}

	return
}

func Parse(version versions.SignedVersion, permissions string) (sp SignedPermissions) {
	sp = SignedPermissions{
		hasValues:   false,
		permissions: [numPermissions]string{},

		SignedVersion: version,
	}

	spMap := signedPermissionMap()
	splitPermissions := strings.Split(permissions, "")
	for _, permission := range splitPermissions {
		conformedPermission := strings.ToLower(permission)
		if spec, ok := spMap[conformedPermission]; ok {
			sp.permissions[spec.Index] = conformedPermission
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

func signedPermissionMap() map[string]signedPermissionSpec {
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas#permissions-for-a-directory-container-or-blob
	return map[string]signedPermissionSpec{
		"r": {
			OpName:        "Read",
			OpDescription: "Read the content, block list, properties, and metadata of any blob in the container or directory. Use a blob as the source of a copy operation.",
			Index:         0,
			APIVersion:    versions.VAll,
		},
		"a": {
			OpName:        "Add",
			OpDescription: "Add a block to an append blob.",
			Index:         1,
			APIVersion:    versions.VAll,
		},
		"c": {
			OpName:        "Create",
			OpDescription: "Write a new blob, snapshot a blob, or copy a blob to a new blob.",
			Index:         2,
			APIVersion:    versions.VAll,
		},
		"w": {
			OpName:        "Write",
			OpDescription: "Create or write content, properties, metadata, or block list. Snapshot or lease the blob. Resize the blob (page blob only). Use the blob as the destination of a copy operation.",
			Index:         3,
			APIVersion:    versions.VAll,
		},
		"d": {
			OpName:        "Delete",
			OpDescription: "Delete a blob. For version 2017-07-29 and later, the Delete permission also allows breaking a lease on a blob. For more information, see the Lease Blob operation.",
			Index:         4,
			APIVersion:    versions.VAll,
		},
		"x": {
			OpName:        "Delete version",
			OpDescription: "Delete a blob version.",
			Index:         5,
			APIVersion:    versions.V2019_12_12,
		},
		"y": {
			OpName:        "Permanent delete",
			OpDescription: "Permanently delete a blob snapshot or version.",
			Index:         6,
			APIVersion:    versions.V2020_02_10,
		},
		"l": {
			OpName:        "List",
			OpDescription: "List blobs non-recursively.",
			Index:         7,
			APIVersion:    versions.VAll,
		},
		"t": {
			OpName:        "Tags",
			OpDescription: "Read or write the tags on a blob.",
			Index:         8,
			APIVersion:    versions.V2019_12_12,
		},
		"m": {
			OpName:        "Move",
			OpDescription: "Move a blob or a directory and its contents to a new location. This operation can optionally be restricted to the owner of the child blob, directory, or parent directory if the `saoid` parameter is included on the SAS token and the sticky bit is set on the parent directory.",
			Index:         9,
			APIVersion:    versions.V2020_02_10,
		},
		"e": {
			OpName:        "Execute",
			OpDescription: "Get the system properties and, if the hierarchical namespace is enabled for the storage account, get the POSIX ACL of a blob. If the hierarchical namespace is enabled and the caller is the owner of a blob, this permission grants the ability to set the owning group, POSIX permissions, and POSIX ACL of the blob. Does not permit the caller to read user-defined metadata.",
			Index:         10,
			APIVersion:    versions.V2020_02_10,
		},
		"o": {
			OpName:        "Ownership",
			OpDescription: "When the hierarchical namespace is enabled, this permission enables the caller to set the owner or the owning group, or to act as the owner when renaming or deleting a directory or blob within a directory that has the sticky bit set.",
			Index:         11,
			APIVersion:    versions.V2020_02_10,
		},
		"p": {
			OpName:        "Permissions",
			OpDescription: "When the hierarchical namespace is enabled, this permission allows the caller to set permissions and POSIX ACLs on directories and blobs.",
			Index:         12,
			APIVersion:    versions.V2020_02_10,
		},
	}
}
