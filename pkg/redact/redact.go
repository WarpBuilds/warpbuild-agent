// Package redact removes presigned URL secrets from errors before they are logged.
package redact

import (
	"errors"
	"net/url"
	"strings"
)

// URLError removes the query, which holds the signature and token of a
// presigned URL, from the URL in a *url.Error.
func URLError(err error) error {
	if urlErr, ok := errors.AsType[*url.Error](err); ok {
		if base, _, found := strings.Cut(urlErr.URL, "?"); found {
			urlErr.URL = base + "?REDACTED"
		}
	}
	return err
}
