package tuya

import (
	"bytes"
	"encoding/binary"
)

type Packet interface {
	ReadUint32BE() uint32
	ReadByte(int) []byte
	ReadByteRaw(int, int) []byte
	WriteUInt32BE(i uint32)
	WriteByte([]byte)
	Rewind()
	GetBuffer() []byte
}

type packet struct {
	buffer []byte
	pos    int
}

func NewPacketFromByte(b []byte) Packet {
	return &packet{buffer: b, pos: 0}
}

func NewPacket() Packet {
	return &packet{pos: 0}
}

func (p *packet) ReadUint32BE() uint32 {
	i := binary.BigEndian.Uint32(p.buffer[p.pos : p.pos+4])
	p.pos += 4
	return i
}

func (p *packet) ReadByte(i int) []byte {
	b := p.buffer[p.pos : p.pos+i]
	p.pos += i
	return b
}

func (p *packet) ReadByteRaw(f int, t int) []byte {
	return p.buffer[f:t]
}

func (p *packet) WriteUInt32BE(n uint32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, n)
	b := buf.Bytes()
	for i := 0; i < len(b); i++ {
		p.buffer = append(p.buffer, b[i])
	}
	p.pos += len(b)
}

func (p *packet) WriteByte(b []byte) {
	for i := 0; i < len(b); i++ {
		p.buffer = append(p.buffer, b[i])
	}
	p.pos += len(b)
}

func (p *packet) Rewind() {
	p.pos = 0
}

func (p *packet) GetBuffer() []byte {
	return p.buffer
}
