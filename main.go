package main

import (
	"fmt"
	"log"
)

type Doc struct {
	ID   int32
	Data map[string]interface{}
}

type RevertIndex struct {
	Data map[string]map[interface{}][]int32 // fieldName => termName => positing list
}

type Query func(ri *RevertIndex) ([]int32, error)

func Eq(k string, v interface{}) Query {
	return func(ri *RevertIndex) ([]int32, error) {
		return ri.Data[k][v], nil
	}
}

func NewRI() *RevertIndex {
	return &RevertIndex{
		Data: make(map[string]map[interface{}][]int32),
	}
}

func (ri *RevertIndex) MultiQuery(qs ...Query) ([]int32, error) {
	var mergeDocIDs []int32

	for _, q := range qs {
		docIDs, err := ri.Query(q)
		if err != nil {
			return nil, err
		}

		if mergeDocIDs == nil {
			mergeDocIDs = docIDs
		} else {
			mergeDocIDs = mergeInt32Slice(mergeDocIDs, docIDs)
		}
	}

	return mergeDocIDs, nil
}

func mergeInt32Slice(a []int32, b []int32) []int32 {
	retSlice := make([]int32, 0)
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				retSlice = append(retSlice, a[i])
			}
		}
	}

	return retSlice
}

func (ri *RevertIndex) Query(q Query) ([]int32, error) {
	return q(ri)
}

func (ri *RevertIndex) Index(doc *Doc) error {
	for f, v := range doc.Data {
		field, ok := ri.Data[f]
		if !ok {
			ri.Data[f] = map[interface{}][]int32{v: []int32{doc.ID}}
		} else {
			_, ok := field[v]
			if !ok {
				field[v] = []int32{doc.ID}
			} else {
				field[v] = append(field[v], doc.ID)
			}

		}
	}
	return nil
}

func main() {
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

	fmt.Println(docIDs)
}
