package main

import iradix "github.com/hashicorp/go-immutable-radix"

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
