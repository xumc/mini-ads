package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Converter(t *testing.T) {
	b := Uint32ToBytes(222)
	assert.Equal(t, uint32(222), BytesToUint32(b))
}
