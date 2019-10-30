package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TrieDic_1(t *testing.T) {
	tt := NewTrieDic()
	tt.Add("ZH-BJ", DocID(11111))
	assert.Equal(t, 1, tt.Get("ZH-BJ").Len())
	iter := tt.Get("ZH-BJ").Iterator()
	iter.Next()
	assert.Equal(t, DocID(11111), iter.Key())
}
