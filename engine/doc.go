package engine

type DocID int32

type Doc struct {
	ID   DocID
	Data map[string]interface{}
}
