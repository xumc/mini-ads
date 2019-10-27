package main

type RevertIndex struct {
	fieldMap map[string]Field
	dics     map[string]Dic
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

func (ri *RevertIndex) Query(q Query) ([]int32, error) {
	return q(ri)
}

func (ri *RevertIndex) Index(doc *Doc) error {
	for f, v := range doc.Data {
		ri.AddForField(f).Add(v, doc.ID)
	}
	return nil
}
