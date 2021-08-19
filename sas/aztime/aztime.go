package aztime

import (
	// Standard Library Imports
	"errors"
	"net/url"
	"time"
)

const (
	ParamKeySignedStart  = "st"
	ParamKeySignedExpiry = "se"
)

var (
	ErrDateTimeEmpty = errors.New("datetime provided to parse is empty")
)

// ParseISO8601DateTime provides a much more CLI user-friendly time parser which
// attempts to parse from least-to-greatest precision, failing if it .
func ParseISO8601DateTime(dateTime string) (t time.Time, err error) {
	if dateTime == "" {
		return time.Time{}, ErrDateTimeEmpty
	}

	inputFormats := []string{
		"2006-01-02",             // YYYY-MM-DD
		"2006-01-02Z07:00",       // YYYY-MM-DD<TZDSuffix>
		"2006-01-02T15:04",       // YYYY-MM-DDThh:mm
		"2006-01-02T15:04Z07:00", // YYYY-MM-DDThh:mm<TZDSuffix>
		"2006-01-02T15:04:05",    // YYYY-MM-DDThh:mm:ss
		time.RFC3339,             // YYYY-MM-DDThh:mm:ss<TZDSuffix>
	}

	for _, format := range inputFormats {
		if t, err = time.Parse(format, dateTime); err == nil {
			return t, nil
		}
	}

	return t, err
}

func ToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func GetParam(paramKey string, t time.Time) (timeParam string) {
	if !t.IsZero() {
		params := &url.Values{}
		params.Add(paramKey, ToString(t))

		timeParam = params.Encode()
	}

	return
}

func GetURLDecodedParam(paramKey string, t time.Time) (decodedTimeParam string) {
	if !t.IsZero() {
		decodedTimeParam, _ = url.QueryUnescape(GetParam(paramKey, t))
	}

	return
}
