package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
)

type Codec interface {
	// Encode the given bytes into a protocol.
	Encode(buf []byte) ([]byte, error)
	// Decode the protocol into bytes.
	Decode(buf []byte) ([]byte, error)
	// Unpack the protocol into bytes.
	Unpack(conn gnet.Conn) ([]byte, error)
}

const (
	magic     = 1314
	magicSize = 2
	bodySize  = 4
)

var (
	ErrIncompletePacket = errors.New("incomplete protocol")
	ErrInvalidMagic     = errors.New("invalid magic number")

	magicBytes []byte
)

func init() {
	binary.BigEndian.PutUint16(magicBytes, uint16(magic))
}

type FixedHeaderCodec struct {
}

func (p FixedHeaderCodec) Encode(buf []byte) ([]byte, error) {
	offset := magicSize + bodySize
	msgLen := offset + len(buf)

	data := make([]byte, msgLen)
	copy(data, magicBytes)

	binary.BigEndian.PutUint32(data[magicSize:offset], uint32(len(buf)))
	copy(data[offset:msgLen], buf)
	return data, nil
}

func (p FixedHeaderCodec) Decode(buf []byte) ([]byte, error) {
	offset := magicSize + bodySize
	if len(buf) < offset {
		return nil, ErrIncompletePacket
	}
	if !bytes.Equal(magicBytes, buf[:magicSize]) {
		return nil, ErrInvalidMagic
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicSize:offset])
	msgLen := offset + int(bodyLen)
	if len(buf) < msgLen {
		return nil, ErrIncompletePacket
	}
	return buf[offset:msgLen], nil
}

func (p FixedHeaderCodec) Unpack(c gnet.Conn) ([]byte, error) {
	// 读取固定报头（magic_size + body_size）
	offset := magicSize + bodySize
	buf, _ := c.Peek(offset)
	if len(buf) < offset {
		return nil, ErrIncompletePacket
	}
	if !bytes.Equal(magicBytes, buf[:magicSize]) {
		return nil, ErrInvalidMagic
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicSize:offset])
	msgLen := offset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}

	// 取出数据包，再丢弃掉
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)
	return buf[offset:msgLen], nil
}
