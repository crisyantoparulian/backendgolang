package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	test := 8
	res := ToPointer(test)

	assert.Equal(t, res, &test)
}
