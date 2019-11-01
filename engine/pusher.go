package engine

import (
	"encoding/json"
	"github.com/ryszard/goskiplist/skiplist"
	"io"
	"os"
)

type pusher struct {
	ri *RevertIndex

	dataFile *os.File
}

func NewPusher(ri *RevertIndex) *pusher {
	return &pusher{
		ri:       ri,
		dataFile: nil,
	}
}

func (p *pusher) WriteTo(w io.Writer) (int, error) {
	dicPostingListOffsetMap := make(map[string]map[interface{}]int64)
	for f, dic := range p.ri.dics {
		err := dic.Range(func(k interface{}, postingList *skiplist.Set) error {
			offset, err := p.savePostingList(postingList)
			if err != nil {
				return err
			}
			_, ok := dicPostingListOffsetMap[f]
			if !ok {
				dicPostingListOffsetMap[f] = make(map[interface{}]int64)
			}
			dicPostingListOffsetMap[f][k] = offset
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	fieldDicOffsetMap := make(map[string]int64)
	for f, dic := range dicPostingListOffsetMap {
		dicOffset, err := p.saveDic(dic)
		if err != nil {
			return 0, err
		}

		fieldDicOffsetMap[f] = dicOffset
	}

	dicsDataOffset, err := p.saveDicsData(fieldDicOffsetMap)
	if err != nil {
		return 0, err
	}

	fieldMapOffset, err := p.saveFieldMap()
	if err != nil {
		return 0, err
	}

	err = p.saveMetadata(&Metadata{
		DicsOffset:     dicsDataOffset,
		FieldMapOffset: fieldMapOffset,
	})
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (p *pusher) saveFieldMap() (int64, error) {
	info, err := p.dataFile.Stat()
	if err != nil {
		return 0, err
	}

	retOffset := info.Size()

	p.dataFile.Write(Uint32ToBytes(uint32(len(p.ri.fieldMap))))

	for f, s := range p.ri.fieldMap {
		kbytes := []byte(f)
		sbytes, err := json.Marshal(s)
		if err != nil {
			return 0, err
		}

		keySize := uint32(len(kbytes))
		valueSize := uint32(len(sbytes))

		line := append(Uint32ToBytes(keySize), Uint32ToBytes(valueSize)...)
		line = append(line, kbytes...)
		line = append(line, sbytes...)

		_, err = p.dataFile.Write(line)
		if err != nil {
			return 0, err
		}
	}
	return retOffset, nil
}

func (p *pusher) WriteToFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	p.dataFile = file

	_, err = p.WriteTo(p.dataFile)
	return err
}

func (p *pusher) savePostingList(positings *skiplist.Set) (int64, error) {
	info, err := p.dataFile.Stat()
	if err != nil {
		return 0, err
	}

	offset := info.Size()

	count := positings.Len()
	bytes := make([]byte, 0, 4+count)
	bytes = append(bytes, Uint32ToBytes(uint32(count))...)

	unboundIterator := positings.Iterator()
	for unboundIterator.Next() {
		docID := unboundIterator.Key().(DocID)
		bytes = append(bytes, Int32ToBytes(int32(docID))...)
	}

	_, err = p.dataFile.Write(bytes)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (p *pusher) saveDic(dic map[interface{}]int64) (int64, error) {
	info, err := p.dataFile.Stat()
	if err != nil {
		return 0, err
	}
	retOffset := info.Size()

	p.dataFile.Write(Uint32ToBytes(uint32(len(dic))))

	for k, offset := range dic {
		var keyBytes []byte
		switch v := k.(type) {
		case string:
			keyBytes = []byte(v)
		case int:
			keyBytes = IntToBytes(v)
		}

		offsetBytes := Int64ToBytes(offset)

		keySize := uint32(len(keyBytes))
		offsetSize := uint32(len(offsetBytes))

		line := append(Uint32ToBytes(keySize), Uint32ToBytes(offsetSize)...)
		line = append(line, keyBytes...)
		line = append(line, offsetBytes...)

		_, err := p.dataFile.Write(line)
		if err != nil {
			return 0, err
		}
	}
	return retOffset, nil
}

func (p *pusher) saveDicsData(fieldDics map[string]int64) (int64, error) {
	info, err := p.dataFile.Stat()
	if err != nil {
		return 0, err
	}
	retOffset := info.Size()

	p.dataFile.Write(Uint32ToBytes(uint32(len(fieldDics))))

	for f, offset := range fieldDics {
		kbytes := []byte(f)
		sbytes := Int64ToBytes(offset)

		keySize := uint32(len(kbytes))
		valueSize := uint32(len(sbytes))

		line := append(Uint32ToBytes(keySize), Uint32ToBytes(valueSize)...)
		line = append(line, kbytes...)
		line = append(line, sbytes...)

		_, err := p.dataFile.Write(line)
		if err != nil {
			return 0, err
		}
	}

	return retOffset, nil
}

func (p *pusher) saveMetadata(md *Metadata) error {
	_, err := p.dataFile.Write(md.toBytes())
	if err != nil {
		return err
	}

	return nil
}
