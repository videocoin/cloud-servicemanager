package service

import (
	"regexp"

	"github.com/videocoin/cloud-pkg/api/resources"
)

// ErrInvalidName indicates that the service name is invalid.
var ErrInvalidName = resources.PatternError(NamePattern.String())

// NamePattern is the service name pattern.
var NamePattern = regexp.MustCompile(`^[a-z]{3,20}\.videocoin\.network$`)

// IsValidName whether a service name is valid.
func IsValidName(name string) bool {
	return NamePattern.MatchString(name)
}
