package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TrieDic_1(t *testing.T) {
	tt := NewTrieDic()
	tt.Add("ZH-BJ", int32(11111))
	assert.Equal(t, []int32{11111}, tt.Get("ZH-BJ"))
}
