package buffer

import "sync"

const DefaultSmallBytesLen = 4

type SmallBytes struct {
	Data []byte
}

func (s *SmallBytes) reset() {
	s.Data = s.Data[:0]
}

func (s *SmallBytes) Release() {
	s.reset()
	SmallBytesPool.Put(s)
}

func NewSmallBytes(len int64) *SmallBytes {
	return &SmallBytes{
		Data: make([]byte, len),
	}
}

var SmallBytesPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return NewSmallBytes(DefaultSmallBytesLen)
	},
}
