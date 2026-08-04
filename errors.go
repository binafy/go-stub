package stub

import (
	"errors"
	"strings"
)

// ErrExists is returned by Generate and the other file-writing functions when
// the destination already exists and no policy option (WithForce,
// WithSkipExisting, WithAppend) was supplied.
var ErrExists = errors.New("stub: destination already exists")

// ErrMissingKeys is the sentinel that a *MissingKeysError unwraps to, so
// callers can test errors.Is(err, ErrMissingKeys) without depending on the
// concrete type.
var ErrMissingKeys = errors.New("stub: unresolved placeholders")

// ErrUnsafePath is returned by GenerateDir and GenerateDirFS when a rendered
// file name would place the output outside the destination directory (for
// example because a placeholder value contains "../"). It is a guard against
// path traversal from untrusted replacement values.
var ErrUnsafePath = errors.New("stub: generated path escapes the destination directory")

// MissingKeysError reports the placeholder keys that were left unresolved
// during a render performed with WithStrict. Keys are listed in first-seen
// order. Retrieve them with errors.As.
type MissingKeysError struct {
	Keys []string
}

// Error implements the error interface.
func (e *MissingKeysError) Error() string {
	return "stub: unresolved placeholders: " + strings.Join(e.Keys, ", ")
}

// Unwrap lets errors.Is(err, ErrMissingKeys) match a *MissingKeysError.
func (e *MissingKeysError) Unwrap() error {
	return ErrMissingKeys
}
