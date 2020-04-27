package context

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroup(t *testing.T) {

	r := &router{}
	g := r.Group("test")
	gg := g.Group("//gogo")

	assert.Equal(t, "test/gogo/:id", gg.JoinedPath("/:id"))
}
