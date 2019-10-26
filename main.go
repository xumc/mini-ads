package main

import (
	"errors"
	"github.com/hashicorp/go-immutable-radix"
)

type Doc struct {
	ID   int32
	Data map[string]interface{}
}

type FieldType int

const (
	_ FieldType = iota
	Number
	Bytes
	Bool
)

type DicType int

const (
	_ DicType = iota
	Hash
	Trie
)

type Field struct {
	Name    string
	Type    FieldType
	DicType DicType
}

type Schema []Field

type RevertIndex struct {
	fieldMap map[string]Field
	dics     map[string]Dic
}

type Dic interface {
	Get(key interface{}) []int32
	Add(key interface{}, value int32)
}

type hashDic struct {
	data map[interface{}][]int32
}

func NewHashDic() Dic {
	return &hashDic{
		data: make(map[interface{}][]int32),
	}
}

func (hd *hashDic) Get(key interface{}) []int32 {
	return hd.data[key]
}

func (hd *hashDic) Add(key interface{}, value int32) {
	_, ok := hd.data[key]
	if !ok {
		hd.data[key] = []int32{value,}
		return
	}

	hd.data[key] = append(hd.data[key], value)
}

type trieDic struct {
	data *iradix.Tree
}

func NewTrieDic() Dic {
	return &trieDic{
		data: iradix.New(),
	}
}

func (td *trieDic) Get(key interface{}) []int32 {
	ps, ok := td.data.Get([]byte(key.(string)))
	if !ok {
		return []int32{}
	}
	return ps.([]int32)
}

func (td *trieDic) Add(key interface{}, value int32) {
	ps := td.Get(key)
	ps = append(ps, value)

	td.data, _, _ = td.data.Insert([]byte(key.(string)), ps)
}

func (ri *RevertIndex) GetDic(field string) Dic {
	return ri.dics[field]
}

func (ri *RevertIndex) AddForField(field string) Dic {
	fs, ok := ri.dics[field]
	if !ok {
		switch ri.fieldMap[field].DicType {
		case Hash:
			ri.dics[field] = NewHashDic()
		case Trie:
			ri.dics[field] = NewTrieDic()
		}
		fs = ri.dics[field]
	}

	return fs
}

type Query func(ri *RevertIndex) ([]int32, error)

func Eq(k string, v interface{}) Query {
	return func(ri *RevertIndex) ([]int32, error) {
		return ri.GetDic(k).Get(v), nil
	}
}

func StartWith(k string, v interface{}) Query {
	return func(ri *RevertIndex) ([]int32, error) {
		switch vv := ri.GetDic(k).(type) {
		case *hashDic:
			return nil, errors.New("unsuport startWith Query in hash type field.")
		case *trieDic:
			retDocIDs := make([]int32, 0)
			vv.data.Root().WalkPrefix([]byte(v.(string)), func(kkk []byte, vvv interface{}) bool {
				retDocIDs = append(retDocIDs, vvv.([]int32)...)
				return false
			})
			return retDocIDs, nil
		default:
			return nil, errors.New("unsupport field dic type.")
		}
	}
}

func NewRI(schema Schema) *RevertIndex {
	fieldMap := make(map[string]Field)
	for _, f := range schema {
		fieldMap[f.Name] = f
	}

	return &RevertIndex{
		fieldMap: fieldMap,
		dics:     make(map[string]Dic),
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
			mergeDocIDs = intersectInt32Slice(mergeDocIDs, docIDs)
		}
	}

	return mergeDocIDs, nil
}

func intersectInt32Slice(a []int32, b []int32) []int32 {
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
		ri.AddForField(f).Add(v, doc.ID)
	}
	return nil
}

func main() {
}
