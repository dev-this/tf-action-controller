package webhook

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewEvent asserts arguments from the constructor call are passed into field values correctly
func TestNewEvent(t *testing.T) {
	want := Event{
		sha:   "a1b2c3fffff",
		ref:   "tags/v3",
		owner: "octocat",
		repo:  "not-a-real-repo",
	}
	got := NewEvent("a1b2c3fffff", "tags/v3", "octocat", "not-a-real-repo")

	assert.EqualValues(t, want, got)
}
