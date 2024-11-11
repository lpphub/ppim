package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
)

const (
	magicNum = 1314
	magicLen = 2
	bodyLen  = 4
)

var (
	ErrInvalidMagic = errors.New("invalid magic number")
	magicBytes      []byte
)

func init() {
	magicBytes = make([]byte, magicLen)
	binary.BigEndian.PutUint16(magicBytes, uint16(magicNum))
}

type FixedFrameCodec struct {
}

func (p FixedFrameCodec) Encode(buf []byte) ([]byte, error) {
	offset := magicLen + bodyLen
	msgLen := offset + len(buf)

	data := make([]byte, msgLen)
	copy(data, magicBytes)

	binary.BigEndian.PutUint32(data[magicLen:offset], uint32(len(buf)))
	copy(data[offset:msgLen], buf)
	return data, nil
}

func (p FixedFrameCodec) Decode(c gnet.Conn) ([]byte, error) {
	offset := magicLen + bodyLen
	buf, _ := c.Peek(offset)
	if len(buf) < offset {
		return nil, ErrIncompletePacket
	}
	if !bytes.Equal(magicBytes, buf[:magicLen]) {
		// 丢弃非法数据包
		_, _ = c.Discard(magicLen)
		return nil, ErrInvalidMagic
	}

	payloadLen := binary.BigEndian.Uint32(buf[magicLen:offset])
	msgLen := offset + int(payloadLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}

	buf, _ = c.Next(msgLen)
	return buf[offset:msgLen], nil
}

func (p FixedFrameCodec) Unpack(buf []byte) ([]byte, error) {
	offset := magicLen + bodyLen
	if len(buf) < offset {
		return nil, ErrIncompletePacket
	}
	if !bytes.Equal(magicBytes, buf[:magicLen]) {
		return nil, ErrInvalidMagic
	}

	payloadLen := binary.BigEndian.Uint32(buf[magicLen:offset])
	msgLen := offset + int(payloadLen)
	if len(buf) < msgLen {
		return nil, ErrIncompletePacket
	}
	return buf[offset:msgLen], nil
}
