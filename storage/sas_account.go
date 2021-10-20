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

package storage

import (
	// Standard Library Imports
	"encoding/base64"
	"net/url"
	"time"

	// Internal Imports
	"github.com/matthewhartstonge/sassy/storage/aztime"
	"github.com/matthewhartstonge/sassy/storage/crypto"
	"github.com/matthewhartstonge/sassy/storage/ips"
	"github.com/matthewhartstonge/sassy/storage/permissions"
	"github.com/matthewhartstonge/sassy/storage/protocols"
	"github.com/matthewhartstonge/sassy/storage/resourcetypes"
	"github.com/matthewhartstonge/sassy/storage/services"
	"github.com/matthewhartstonge/sassy/storage/versions"
)

// NewAccountSAS provides a way to generate an account based Shared Access
// Signature (SAS) token.
func NewAccountSAS(
	storageAccountName string,
	storageAccountKey string,
	signedVersion string,
	signedServices string,
	signedResourceTypes string,
	signedPermissions string,
	signedExpiry string,
	opts ...AccountSASOption,
) (
	accountSAS *AccountSAS,
	err error,
) {
	storageKeyBytes, err := base64.StdEncoding.DecodeString(storageAccountKey)
	if err != nil {
		return nil, ErrDecodingStorageAccountKey
	}

	sv, ok := versions.Parse(signedVersion)
	if !ok {
		return nil, ErrInvalidVersion
	}

	se, err := aztime.ParseISO8601DateTime(signedExpiry)
	if err != nil {
		switch err {
		case aztime.ErrDateTimeEmpty:
			// Bubble up internal known errors.
			return nil, err

		default:
			// Overwrite time parsing errors as an invalid format error.
			return nil, ErrInvalidExpiryDateFormat
		}
	}

	accountSAS = &AccountSAS{
		storageAccountName:  storageAccountName,
		storageAccountKey:   storageKeyBytes,
		SignedVersion:       sv,
		SignedServices:      services.Parse(signedServices),
		SignedResourceTypes: resourcetypes.Parse(signedResourceTypes),
		SignedPermission:    permissions.Parse(sv, signedPermissions),
		SignedExpiry:        se,
	}

	// Inject optional fields
	for _, opt := range opts {
		if err := opt(accountSAS); err != nil {
			return nil, err
		}
	}

	return accountSAS, nil
}

type AccountSASOption func(options *AccountSAS) error

func WithAPIVersion(apiVersion string) AccountSASOption {
	return func(options *AccountSAS) error {
		options.APIVersion = apiVersion

		return nil
	}
}

func WithSignedStart(startDateTime string) AccountSASOption {
	return func(options *AccountSAS) error {
		st, err := aztime.ParseISO8601DateTime(startDateTime)
		if err != nil {
			switch err {
			case aztime.ErrDateTimeEmpty:
				// Bubble up internal known errors.
				return err

			default:
				// Overwrite time parsing errors as an invalid format error.
				return ErrInvalidStartDateFormat
			}
		}

		options.SignedStart = st

		return nil
	}
}

func WithSignedIP(ip string) AccountSASOption {
	return func(options *AccountSAS) error {
		sip, ok := ips.Parse(ip)
		if !ok {
			return ErrInvalidIPv4Format
		}

		options.SignedIP = sip
		return nil
	}
}

func WithSignedProtocols(signedProtocols string) AccountSASOption {
	return func(options *AccountSAS) error {
		options.SignedProtocol = protocols.Parse(signedProtocols)

		return nil
	}
}

type AccountSAS struct {
	storageAccountName  string
	storageAccountKey   []byte
	APIVersion          string
	SignedVersion       versions.SignedVersion
	SignedServices      services.SignedServices
	SignedResourceTypes resourcetypes.SignedResourceTypes
	SignedPermission    permissions.SignedPermissions
	SignedStart         time.Time
	SignedExpiry        time.Time
	SignedIP            ips.SignedIP
	SignedProtocol      protocols.SignedProtocols
}

// Token generates and signs an account based storage SAS token based on the
// stored configuration.
func (o AccountSAS) Token() string {
	params := &url.Values{}
	if o.APIVersion != "" {
		params.Add("api-version", o.APIVersion)
	}

	o.SignedVersion.SetParam(params)
	o.SignedServices.SetParam(params)
	o.SignedResourceTypes.SetParam(params)
	o.SignedPermission.SetParam(params)

	if !o.SignedStart.IsZero() {
		params.Set(
			aztime.ParamKeySignedStart,
			aztime.ToString(o.SignedStart),
		)
	}

	if !o.SignedExpiry.IsZero() {
		params.Set(
			aztime.ParamKeySignedExpiry,
			aztime.ToString(o.SignedExpiry),
		)
	}

	o.SignedIP.SetParam(params)
	o.SignedProtocol.SetParam(params)
	o.signPayload(params)

	return params.Encode()
}

// signPayload generates the required HMAC-SHA256 signature and binds it into
// the provided url params.
func (o AccountSAS) signPayload(params *url.Values) {
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-account-sas#constructing-the-signature-string
	// To construct the signature string for an account SAS, first construct the
	// string-to-sign from the fields comprising the request, then encode the
	// string as UTF-8 and compute the signature using the HMAC-SHA256
	// algorithm.
	//
	// Note:
	// - Fields included in the string-to-sign must be UTF-8, URL-decoded.
	//   - Go by default uses utf-8 encoded strings.
	//   - The `String()` methods ensure no URL encoding is taking place.
	stringToSign := o.storageAccountName + "\n" +
		o.SignedPermission.String() + "\n" +
		o.SignedServices.String() + "\n" +
		o.SignedResourceTypes.String() + "\n" +
		aztime.ToString(o.SignedStart) + "\n" +
		aztime.ToString(o.SignedExpiry) + "\n" +
		o.SignedIP.String() + "\n" +
		o.SignedProtocol.String() + "\n" +
		o.SignedVersion.String() + "\n"

	// Compute HMAC-S256 signature
	signature := crypto.HMACSHA256(
		o.storageAccountKey,
		[]byte(stringToSign),
	)

	params.Add("sig", signature)
}
