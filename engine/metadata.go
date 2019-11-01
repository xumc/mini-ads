package engine

import (
	"unsafe"
)

type Metadata struct {
	DicsOffset     int64
	FieldMapOffset int64
	DocsOffset     int64
	ptr            uintptr
}

// maxAllocSize is the size used when creating array pointers.
const maxAllocSize = 0x7FFFFFF

func NewMetadataFromBytes(bytes []byte) *Metadata {
	m := Metadata{}
	b := ((*[maxAllocSize]byte)(unsafe.Pointer(&m)))[:unsafe.Sizeof(m)]
	copy(b, bytes)
	return &m
}

func (md *Metadata) toBytes() []byte {
	return (*[maxAllocSize]byte)(unsafe.Pointer(md))[:unsafe.Sizeof(*md)]
}
