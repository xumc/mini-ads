package engine

import "errors"

type Query func(ri *RevertIndex) ([]int32, error)

func Eq(field string, key interface{}) Query {
	return func(ri *RevertIndex) ([]int32, error) {
		r := ri.GetDic(field).Get(key)
		return r, nil
	}
}

func StartWith(field string, key interface{}) Query {
	return func(ri *RevertIndex) ([]int32, error) {
		switch vv := ri.GetDic(field).(type) {
		case *hashDic:
			return nil, errors.New("unsuport startWith Query in hash type field.")
		case *trieDic:
			retDocIDs := make([]int32, 0)
			vv.data.Root().WalkPrefix([]byte(key.(string)), func(kkk []byte, vvv interface{}) bool {
				retDocIDs = append(retDocIDs, vvv.([]int32)...)
				return false
			})
			return retDocIDs, nil
		default:
			return nil, errors.New("unsupport field dic type.")
		}
	}
}
