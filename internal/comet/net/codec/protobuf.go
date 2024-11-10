package codec

import (
	"encoding/binary"
	"github.com/panjf2000/gnet/v2"
)

type ProtobufCodec struct {
}

const (
	_pbHeaderLen = 4
)

func (p ProtobufCodec) Encode(buf []byte) ([]byte, error) {
	out := make([]byte, _pbHeaderLen+len(buf))
	binary.BigEndian.PutUint32(out, uint32(len(buf)))
	copy(out[_pbHeaderLen:], buf)
	return out, nil
}

func (p ProtobufCodec) Decode(c gnet.Conn) ([]byte, error) {
	buf, _ := c.Peek(_pbHeaderLen)
	if len(buf) < _pbHeaderLen {
		return nil, ErrIncompletePacket
	}

	payloadLen := binary.BigEndian.Uint32(buf)
	msgLen := _pbHeaderLen + int(payloadLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}

	buf, _ = c.Next(msgLen)
	return buf[_pbHeaderLen:], nil
}

func (p ProtobufCodec) Unpack(buf []byte) ([]byte, error) {
	payloadLen := binary.BigEndian.Uint32(buf[:_pbHeaderLen])
	return buf[_pbHeaderLen:payloadLen], nil
}
