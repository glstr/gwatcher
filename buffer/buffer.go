package buffer

import "sync"

const (
	MaxPacketLen = 1500
)

type PacketBuffer struct {
	Buffer   []byte
	len      int64
	capacity int64
}

func NewPacketBuffer() *PacketBuffer {
	return &PacketBuffer{
		Buffer:   make([]byte, MaxPacketLen),
		len:      0,
		capacity: MaxPacketLen,
	}
}

func (p *PacketBuffer) reset() {
	p.Buffer = p.Buffer[:0]
	p.len = 0
}

func (p *PacketBuffer) Release() {
	p.reset()
	packetBufferPool.Put(p)
}

var packetBufferPool = sync.Pool{
	New: func() interface{} {
		return NewPacketBuffer()
	},
}

func GetPacketBuffer() *PacketBuffer {
	get := packetBufferPool.Get()
	return get.(*PacketBuffer)
}
