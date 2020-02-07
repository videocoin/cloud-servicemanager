package service

import (
	"regexp"

	"github.com/videocoin/common/api/resource"
)

// ErrInvalidName indicates that the service name is invalid.
var ErrInvalidName = resource.PatternError(NamePattern.String())

// NamePattern is the service name pattern.
var NamePattern = regexp.MustCompile(`^[a-z]{3,20}\.videocoin\.network$`)

// IsValidName whether a service name is valid.
func IsValidName(name string) bool {
	return NamePattern.MatchString(name)
}
