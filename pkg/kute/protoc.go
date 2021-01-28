package kute

import (
	"encoding/binary"
)

type Header []byte

const HeaderSize = 32 * 1024

type MessageType uint16

const (
	PASS MessageType = iota
	REPL
)

func (h Header) Encode(mtype MessageType, id int32, payload []byte) {
	binary.BigEndian.PutUint16(h[0:2], uint16(mtype))
	binary.BigEndian.PutUint32(h[6:10], uint32(id))
	for i, b := range payload {
		h[10+i] = b
	}
	binary.BigEndian.PutUint32(h[2:6], uint32(len(payload)))
}

func (h Header) MessageType() MessageType {
	value := binary.BigEndian.Uint16(h[0:2])
	return MessageType(value)
}

func (h Header) Len() uint32 {
	return binary.BigEndian.Uint32(h[2:6])
}

func (h Header) ID() uint32 {
	return binary.BigEndian.Uint32(h[6:10])
}
func (h Header) Payload() []byte {
	return h[10:]
}
