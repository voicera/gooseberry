package urn

import (
	"encoding/json"
	"testing"

	"github.com/voicera/tester/assert"
)

func TestJSONMarshalling(t *testing.T) {
	u := NewURN("NID:NID", "NSS")
	buffer, err := json.Marshal(u)
	if assert.For(t).ThatActual(err).IsNil().Passed() {
		var unmarshalled *URN
		err = json.Unmarshal(buffer, &unmarshalled)
		if assert.For(t).ThatActual(err).IsNil().Passed() {
			assert.For(t).ThatActual(unmarshalled).Equals(u)
			assert.For(t).ThatActualString(unmarshalled.GetNamespaceID()).Equals("NID:NID")
			assert.For(t).ThatActualString(unmarshalled.GetNamespaceSpecificString()).Equals("NSS")
		}
	}
}
