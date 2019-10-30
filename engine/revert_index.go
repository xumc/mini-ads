package engine

import (
	"encoding/json"
	"github.com/ryszard/goskiplist/skiplist"
	"io"
	"os"
)

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

func (ri *RevertIndex) MultiQuery(qs ...Query) ([]DocID, error) {
	var mergedDocIDs = newDocIDSkipList()

	for i, q := range qs {
		slDocIDs, err := ri.query(q)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			mergedDocIDs = slDocIDs
		} else {
			mergedDocIDs = intersectSkipList(mergedDocIDs, slDocIDs)
			if mergedDocIDs.Len() == 0 {
				return []DocID{}, nil
			}
		}
	}

	docIDs := make([]DocID, 0, mergedDocIDs.Len())
	iter := mergedDocIDs.Iterator()
	for iter.Next() {
		docIDs = append(docIDs, iter.Key().(DocID))
	}
	return docIDs, nil
}

func (ri *RevertIndex) query(q Query) (*skiplist.Set, error) {
	sl, err := q(ri)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (ri *RevertIndex) Index(doc *Doc) error {
	for f, v := range doc.Data {
		ri.AddForField(f).Add(v, doc.ID)
	}
	return nil
}

func (ri *RevertIndex) ReLoadFromFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	fieldMapOffset := int64(214)
	err = ri.readFieldMap(file, fieldMapOffset)
	if err != nil {
		return err
	}

	dicsOffset := int64(172)
	err = ri.readDics(file, dicsOffset)
	if err != nil {
		return err
	}

	return nil
}

func (ri *RevertIndex) readFieldMap(file *os.File, fieldMapOffset int64) error {
	file.Seek(fieldMapOffset, 0)

	var lineCountBytes = make([]byte, 4)
	_, err := file.Read(lineCountBytes)
	if err != nil {
		return err
	}
	lineCount := BytesToUint32(lineCountBytes)

	for i := lineCount; i > 0; i-- {
		sizes := make([]byte, 8)
		_, err := file.Read(sizes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		keySize := BytesToUint32(sizes[:4])
		valueSize := BytesToUint32(sizes[4:])

		KVSize := keySize + valueSize

		line := make([]byte, KVSize)
		_, err = file.Read(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		key := string(line[:keySize])
		field := Field{}
		err = json.Unmarshal(line[keySize:], &field)
		if err != nil {
			return err
		}

		ri.fieldMap[key] = field

	}
	return nil
}

func (ri *RevertIndex) readDics(file *os.File, dicsOffset int64) error {
	file.Seek(dicsOffset, 0)

	var lineCountBytes = make([]byte, 4)
	_, err := file.Read(lineCountBytes)
	if err != nil {
		return err
	}
	lineCount := BytesToUint32(lineCountBytes)

	m := make(map[string]int64)

	for i := lineCount; i > 0; i-- {
		sizes := make([]byte, 8)
		_, err := file.Read(sizes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		keySize := BytesToUint32(sizes[:4])
		valueSize := BytesToUint32(sizes[4:])

		KVSize := keySize + valueSize

		line := make([]byte, KVSize)
		_, err = file.Read(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		field := string(line[:keySize])
		dicOffset := BytesToInt64(line[keySize:])

		m[field] = dicOffset
	}

	for f, offset := range m {
		dic, err := ri.readDic(file, f, offset)
		if err != nil {
			return err
		}

		ri.dics[f] = dic
	}

	return nil
}

func (ri *RevertIndex) readDic(file *os.File, field string, dicOffset int64) (Dic, error) {
	var dic Dic
	switch ri.fieldMap[field].DicType {
	case Hash:
		dic = NewHashDic()
	case Trie:
		dic = NewTrieDic()
	}

	file.Seek(dicOffset, 0)
	var lineCountBytes = make([]byte, 4)
	_, err := file.Read(lineCountBytes)
	if err != nil {
		return nil, err
	}
	lineCount := BytesToUint32(lineCountBytes)

	m := make(map[interface{}]int64)

	for i := lineCount; i > 0; i-- {
		sizes := make([]byte, 8)
		_, err := file.Read(sizes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		keySize := BytesToUint32(sizes[:4])
		valueSize := BytesToUint32(sizes[4:])

		KVSize := keySize + valueSize

		line := make([]byte, KVSize)
		_, err = file.Read(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		key := line[:keySize]
		positingListOffset := BytesToInt64(line[keySize:])

		switch ri.fieldMap[field].Type {
		case Number:
			m[BytesToInt(key)] = positingListOffset
		case Bytes:
			m[string(key)] = positingListOffset
		case Bool:
			// TODO
		}
	}

	for key, psOffset := range m {
		positingList, err := ri.readPositingList(file, psOffset)
		if err != nil {
			return nil, err
		}

		dic.Set(key, positingList)
	}

	return dic, nil
}

func (ri *RevertIndex) readPositingList(file *os.File, psOffset int64) (*skiplist.Set, error) {
	file.Seek(psOffset, 0)

	countBytes := make([]byte, 4)
	_, err := file.Read(countBytes)
	if err != nil {
		return nil, err
	}
	count := BytesToUint32(countBytes)

	positingListBytes := make([]byte, 4*count)
	_, err = file.Read(positingListBytes)
	if err != nil {
		return nil, err
	}

	positingList := newDocIDSkipList()
	for i := 0; i < int(count); i++ {
		docID := DocID(BytesToInt32(positingListBytes[(i * 4):((i + 1) * 4) ]))
		positingList.Add(docID)
	}

	return positingList, nil
}
