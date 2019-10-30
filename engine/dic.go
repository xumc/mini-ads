package engine

import (
	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/ryszard/goskiplist/skiplist"
)

type Dic interface {
	Get(key interface{}) *skiplist.Set
	Add(key interface{}, value DocID)
	Set(key interface{}, positingList *skiplist.Set)
	Range(func(k interface{}, postingList *skiplist.Set) error) error
}

type hashDic struct {
	data map[interface{}]*skiplist.Set
}

func NewHashDic() Dic {
	return &hashDic{
		data: make(map[interface{}]*skiplist.Set),
	}
}

func (hd *hashDic) Get(key interface{}) *skiplist.Set {
	return hd.data[key]
}

func (hd *hashDic) Add(key interface{}, value DocID) {
	_, ok := hd.data[key]
	if !ok {
		hd.data[key] = newDocIDSkiplist()
	}

	hd.data[key].Add(value)
}

func (hd *hashDic) Set(key interface{}, positingList *skiplist.Set) {
	_, ok := hd.data[key]
	if !ok {
		hd.data[key] = positingList
		return
	}

	hd.data[key] = positingList
}

func (hd *hashDic) Range(fn func(k interface{}, postingList *skiplist.Set) error) error {
	for k, v := range hd.data {
		err := fn(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}


type trieDic struct {
	data *iradix.Tree
}

func NewTrieDic() Dic {
	return &trieDic{
		data: iradix.New(),
	}
}

func (td *trieDic) Get(key interface{}) *skiplist.Set {
	ps, ok := td.data.Get([]byte(key.(string)))
	if !ok {
		return newDocIDSkiplist()
	}
	return ps.(*skiplist.Set)
}

func (td *trieDic) Add(key interface{}, value DocID) {
	ps := td.Get(key)
	ps.Add(value)

	td.data, _, _ = td.data.Insert([]byte(key.(string)), ps)
}

func (td *trieDic) Set(key interface{}, positingList *skiplist.Set) {
	td.data, _, _ = td.data.Insert([]byte(key.(string)), positingList)
}

func (td *trieDic) Range(fn func(k interface{}, postingList *skiplist.Set) error) error {
	it := td.data.Root().Iterator()

	for key, val, ok := it.Next(); ok; key, val, ok = it.Next() {
		err := fn(string(key), val.(*skiplist.Set))
		if err != nil {
			return err
		}
	}

	return nil
}

func newDocIDSkiplist() *skiplist.Set {
	return skiplist.NewCustomSet(func(l, r interface{}) bool {
		return int32(l.(DocID)) < int32(r.(DocID))
	})
}
