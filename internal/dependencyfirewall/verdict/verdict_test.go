//go:build !integration

package verdict

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryKey(t *testing.T) {
	e := Entry{Package: "foo", Version: "1.2.3", Verdict: Blocked}
	assert.Equal(t, "foo@1.2.3:blocked", e.Key())
}
