package engine

import (
	"errors"
	"github.com/ryszard/goskiplist/skiplist"
)

type Query func(ri *RevertIndex) (*skiplist.Set, error)

func Eq(field string, key interface{}) Query {
	return func(ri *RevertIndex) (*skiplist.Set, error) {
		r := ri.GetDic(field).Get(key)
		return r, nil
	}
}

func StartWith(field string, key interface{}) Query {
	return func(ri *RevertIndex) (*skiplist.Set, error) {
		switch vv := ri.GetDic(field).(type) {
		case *hashDic:
			return nil, errors.New("unsuport startWith Query in hash type field.")
		case *trieDic:
			retDocIDs := newDocIDSkipList()
			vv.data.Root().WalkPrefix([]byte(key.(string)), func(kkk []byte, vvv interface{}) bool {
				sl := vvv.(*skiplist.Set)
				unboundIterator := sl.Iterator()
				for unboundIterator.Next() {
					retDocIDs.Add(unboundIterator.Key())
				}
				return false
			})
			return retDocIDs, nil
		default:
			return nil, errors.New("unsupport field dic type.")
		}
	}
}
