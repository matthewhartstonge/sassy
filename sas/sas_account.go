package sas

import (
	// Standard Library Imports
	"encoding/base64"
	"net/url"
	"time"

	// Internal Imports
	"github.com/matthewhartstonge/sassy/sas/aztime"
	"github.com/matthewhartstonge/sassy/sas/crypto"
	"github.com/matthewhartstonge/sassy/sas/permissions"
	"github.com/matthewhartstonge/sassy/sas/protocols"
	"github.com/matthewhartstonge/sassy/sas/resourcetypes"
	"github.com/matthewhartstonge/sassy/sas/services"
	"github.com/matthewhartstonge/sassy/sas/versions"
)

const (
	// TODO: package signed-ip
	paramKeySignedIP = "sip"
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
	options *AccountSASOptions,
	err error,
) {
	storageKeyBytes, err := base64.StdEncoding.DecodeString(storageAccountKey)
	if err != nil {
		return nil, ErrInvalidStorageAccountKey
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

	ss := services.Parse(signedServices)
	srt := resourcetypes.Parse(signedResourceTypes)
	sp := permissions.Parse(sv, signedPermissions)

	options = &AccountSASOptions{
		storageAccountName:  storageAccountName,
		storageAccountKey:   storageKeyBytes,
		SignedVersion:       sv,
		SignedServices:      ss,
		SignedResourceTypes: srt,
		SignedPermission:    sp,
		SignedExpiry:        se,
	}

	// Inject user options
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}

	return options, nil
}

type AccountSASOption func(options *AccountSASOptions) error

func WithAPIVersion(apiVersion string) AccountSASOption {
	return func(options *AccountSASOptions) error {
		options.ApiVersion = apiVersion

		return nil
	}
}

func WithSignedStart(startDateTime string) AccountSASOption {
	return func(options *AccountSASOptions) error {
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
	return func(options *AccountSASOptions) error {
		// TODO: parse signed ip using signed ip package
		options.SignedIP = ip

		return nil
	}
}

func WithSignedProtocols(signedProtocols string) AccountSASOption {
	return func(options *AccountSASOptions) error {
		options.SignedProtocol = protocols.Parse(signedProtocols)

		return nil
	}
}

type AccountSASOptions struct {
	storageAccountName  string
	storageAccountKey   []byte
	ApiVersion          string
	SignedVersion       versions.SignedVersion
	SignedServices      services.Services
	SignedResourceTypes resourcetypes.ResourceTypes
	SignedPermission    permissions.SignedPermissions
	SignedStart         time.Time
	SignedExpiry        time.Time
	SignedIP            string // TODO: use package signed-ip type
	SignedProtocol      protocols.Protocols
}

func (o AccountSASOptions) GetToken() string {
	params := &url.Values{}
	if o.ApiVersion != "" {
		params.Add("api-version", o.ApiVersion)
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

	if o.SignedIP != "" {
		// todo: inject parsed signed ip
		params.Set(paramKeySignedIP, o.SignedIP)
	}

	o.SignedProtocol.SetParam(params)
	o.signPayload(params)

	return params.Encode()
}

func (o AccountSASOptions) signPayload(params *url.Values) {
	// Refer: https://docs.microsoft.com/en-us/rest/api/storageservices/create-account-sas#constructing-the-signature-string
	// To construct the signature string for an account SAS, first construct the
	// string-to-sign from the fields comprising the request, then encode the
	// string as UTF-8 and compute the signature using the HMAC-SHA256
	// algorithm.
	//
	// Note:
	// - Fields included in the string-to-sign must be UTF-8, URL-decoded.
	//   - Go by default uses utf-8 encoded strings.
	//   - The `ToString()` methods ensure no URL encoding is taking place.
	stringToSign := o.storageAccountName + "\n" +
		o.SignedPermission.ToString() + "\n" +
		o.SignedServices.ToString() + "\n" +
		o.SignedResourceTypes.ToString() + "\n" +
		aztime.ToString(o.SignedStart) + "\n" +
		aztime.ToString(o.SignedExpiry) + "\n" +
		o.SignedIP + "\n" + // TODO: use packaged signed-ip toString()
		o.SignedProtocol.ToString() + "\n" +
		o.SignedVersion.ToString() + "\n"

	// Compute HMAC-S256 signature
	signature := crypto.HMACSHA256(
		o.storageAccountKey,
		[]byte(stringToSign),
	)

	params.Add("sig", signature)
}
