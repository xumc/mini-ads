package engine

import "github.com/ryszard/goskiplist/skiplist"

func newDocIDSkipList() *skiplist.Set {
	return skiplist.NewCustomSet(func(l, r interface{}) bool {
		return int32(l.(DocID)) < int32(r.(DocID))
	})
}

func intersectSkipList(a *skiplist.Set, b *skiplist.Set) *skiplist.Set {
	retDocIDs := newDocIDSkipList()

	aIter := a.Iterator()
	for aIter.Next() {
		bIter := b.Iterator()
		for bIter.Next() {
			if aIter.Key() == bIter.Key() {
				retDocIDs.Add(aIter.Key())
			}
		}
	}

	return retDocIDs
}
