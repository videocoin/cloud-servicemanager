package key

import (
	"regexp"
	"strings"

	guuid "github.com/google/uuid"
	acc "github.com/videocoin/cloud-pkg/api/resources/serviceaccount"
	cstr "github.com/videocoin/cloud-pkg/strings"

	"github.com/videocoin/cloud-pkg/api/resources"
)

const (
	// CollectionID is the identifier of the resource that contains a list of
	// service account keys.
	CollectionID = "keys"
)

var (
	// ErrInvalidID indicates that the service account key identifier is invalid.
	ErrInvalidID = resources.PatternError(IDPattern.String())
	// ErrInvalidName indicates that the service account key name is invalid.
	ErrInvalidName = resources.PatternError(NamePattern.String())
)

var (
	// IDPattern is the key identifier pattern.
	IDPattern = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	// NamePattern is the key name pattern.
	NamePattern = regexp.MustCompile(`^projects/(([a-z0-9][-a-z0-9]{3,48}[a-z0-9])|\-)/serviceAccounts/[a-z][-a-z0-9]{4,28}[a-z0-9]@[a-z0-9][-a-z0-9]{3,48}[a-z0-9]\.vserviceaccount\.com/keys/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
)

// Name is the key's name.
type Name string

// ID returns the key's identifier.
func (n Name) ID() string {
	return strings.SplitN(string(n), resources.NameSeparator, 6)[5]
}

// ParseName parses a service account name.
func ParseName(name string) (Name, error) {
	if ok := isValidName(name); !ok {
		return "", ErrInvalidName
	}
	return Name(name), nil
}

// NewName returns the key's name given a project identifier and a
// service account email.
func NewName(projID string, accEmail string, keyID string) Name {
	return Name(cstr.JoinWithSeparator(resources.NameSeparator, string(acc.NewName(projID, accEmail)), CollectionID, keyID))
}

// NewNameWildcard returns the service account's key name with a wildcard for
// the project id. Requests using `-` as a wildcard for the project identifier
// will infer the project identifier from the account email.
func NewNameWildcard(accEmail, keyID string) Name {
	return Name(cstr.JoinWithSeparator(resources.NameSeparator, string(acc.NewNameWildcard(accEmail)), CollectionID, keyID))
}

// IsValidID reports whether a key identifier is valid.
// Note: current implementation is faster than regexp.
func IsValidID(ID string) bool {
	_, err := guuid.Parse(ID)
	return err == nil
}

// isValidName reports whether a key name is valid.
func isValidName(name string) bool {
	return NamePattern.MatchString(name)
}
