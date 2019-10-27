package engine

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func Test_Normal_1(t *testing.T) {
	schema := Schema{
		{
			"geo",
			Bytes,
			Trie,
		},
		{
			"age",
			Number,
			Hash,
		},
	}

	ri := NewRI(schema)

	ri.Index(&Doc{
		ID: 1,
		Data: map[string]interface{}{
			"geo": "ZH-BJ",
			"age": 18,
		},
	})

	ri.Index(&Doc{
		ID: 2,
		Data: map[string]interface{}{
			"geo": "ZH-SH",
			"age": 20,
		},
	})

	ri.Index(&Doc{
		ID: 3,
		Data: map[string]interface{}{
			"geo": "ZH-BJ",
			"age": 20,
		},
	})

	ri.Index(&Doc{
		ID: 4,
		Data: map[string]interface{}{
			"geo": "US-NYC",
			"age": 20,
		},
	})

	p := NewPusher(ri)
	defer os.Remove("data")
	p.WriteToFile("data")

	ri = NewRI(schema)
	err := ri.ReLoadFromFile("data")
	assert.Nil(t, err)

	docIDs, err := ri.MultiQuery(Eq("geo", "ZH-BJ"), Eq("age", 20))
	if err != nil {
		log.Fatal(err)
	}

	assert.Len(t, docIDs, 1)
	assert.Equal(t, int32(3), docIDs[0])

	docIDs, err = ri.MultiQuery(StartWith("geo", "ZH"), Eq("age", 20))
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(t, docIDs, 2)
	assert.Equal(t, []int32{3, 2,}, docIDs)
}
