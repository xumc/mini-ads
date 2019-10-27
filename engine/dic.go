package engine

import iradix "github.com/hashicorp/go-immutable-radix"

type Dic interface {
	Get(key interface{}) []int32
	Add(key interface{}, value int32)
	Set(key interface{}, positingList []int32)
	Range(func(k interface{}, postingList []int32) error) error
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

func (hd *hashDic) Set(key interface{}, positingList []int32) {
	_, ok := hd.data[key]
	if !ok {
		hd.data[key] = positingList
		return
	}

	hd.data[key] = positingList
}

func (hd *hashDic) Range(fn func(k interface{}, postingList []int32) error) error {
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

func (td *trieDic) Set(key interface{}, positingList []int32) {
	td.data, _, _ = td.data.Insert([]byte(key.(string)), positingList)
}

func (td *trieDic) Range(fn func(k interface{}, postingList []int32) error) error {
	it := td.data.Root().Iterator()

	for key, val, ok := it.Next(); ok; key, val, ok = it.Next() {
		err := fn(string(key), val.([]int32))
		if err != nil {
			return err
		}
	}

	return nil
}
