package main

import "errors"

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
