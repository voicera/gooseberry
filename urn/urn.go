// Package urn provides a Uniform Resource Name that implemnets RFC8141.
package urn

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	scheme = "urn:"
)

// URN represents a Uniform Resource Name that implemnets RFC8141
// https://www.ietf.org/rfc/rfc8141.txt.
type URN struct {
	namespaceID             string
	namespaceSpecificString string
}

// NewURN creates a new URN with the specified namespace ID
// and namespace-specific string.
func NewURN(namespaceID string, namespaceSpecificString string) *URN {
	return &URN{namespaceID: namespaceID, namespaceSpecificString: namespaceSpecificString}
}

// TryParseString attempts to create a new URN from the specified string.
func TryParseString(urn string) (*URN, bool) {
	parts := strings.Split(urn, ":")
	if len(parts) < 3 || parts[0] != "urn" {
		return nil, false
	}
	return NewURN(parts[1], strings.Join(parts[2:], ":")), true
}

// GetNamespaceID returns the URN's namespace ID (NID).
func (urn *URN) GetNamespaceID() string {
	return urn.namespaceID
}

// GetNamespaceSpecificString returns the URN's namespace-specific string (NSS).
func (urn *URN) GetNamespaceSpecificString() string {
	return urn.namespaceSpecificString
}

func (urn *URN) String() string {
	return scheme + urn.namespaceID + ":" + urn.namespaceSpecificString
}

// MarshalJSON marshals the URN into a JSON string.
func (urn *URN) MarshalJSON() ([]byte, error) {
	return json.Marshal(urn.String())
}

// UnmarshalJSON unmarshals the URN from a JSON string.
func (urn *URN) UnmarshalJSON(b []byte) error {
	var urnString string
	if err := json.Unmarshal(b, &urnString); err != nil {
		return err
	}

	urnString = urnString[len(scheme):] // shift to exclude scheme
	indexOfSeperator := strings.LastIndex(urnString, ":")
	if indexOfSeperator == -1 {
		return errors.New("urn: Unmarshal(malformed URN)")
	}

	urn.namespaceID = urnString[:indexOfSeperator]
	urn.namespaceSpecificString = urnString[indexOfSeperator+1:]
	return nil
}
