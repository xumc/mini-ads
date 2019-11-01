package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_1(t *testing.T) {
	md := &Metadata{DicsOffset: 1, FieldMapOffset: 2, DocsOffset: 3,}
	newMd := &Metadata{}

	newMd = NewMetadataFromBytes(md.toBytes())

	assert.Equal(t, int64(1), newMd.DicsOffset)
	assert.Equal(t, int64(2), newMd.FieldMapOffset)
	assert.Equal(t, int64(3), newMd.DocsOffset)
}
