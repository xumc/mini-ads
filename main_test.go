package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Normal_1(t *testing.T) {
	ri := NewRI()

	ri.Index(&Doc{
		ID: 1,
		Data: map[string]interface{}{
			"geo": "BJ",
			"age": 18,
		},
	})

	ri.Index(&Doc{
		ID: 2,
		Data: map[string]interface{}{
			"geo": "SH",
			"age": 20,
		},
	})

	ri.Index(&Doc{
		ID: 3,
		Data: map[string]interface{}{
			"geo":    "BJ",
			"age":    20,
			"gender": "F",
		},
	})

	docIDs, err := ri.MultiQuery(Eq("geo", "BJ"), Eq("age", 20))
	if err != nil {
		log.Fatal(err)
	}

	assert.Len(t, docIDs, 1)
	assert.Equal(t, 3, docIDs[0])
}
