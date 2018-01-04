package session

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultGet(t *testing.T) {
	s := New("test", nil)
	s.Set("t", "haha")

	assert.Equal(t, "test", s.ID())
	assert.Equal(t, "haha", s.Get("t"))
}
