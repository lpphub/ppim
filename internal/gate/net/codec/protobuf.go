package codec

import (
	"encoding/binary"
	"github.com/panjf2000/gnet/v2"
)

// ProtobufCodec
// * BEFORE (300 bytes)              AFTER (304 bytes)
// * +---------------+               +--------+---------------+
// * | Protobuf Data |-------------->| Length | Protobuf Data |
// * | (300 bytes)   |               | 0x12C  | (300 bytes)   |
// * +---------------+               +--------+---------------+
type ProtobufCodec struct {
}

const (
	_headerLen = 4
)

func (p ProtobufCodec) Encode(buf []byte) ([]byte, error) {
	out := make([]byte, _headerLen+len(buf))
	binary.BigEndian.PutUint32(out, uint32(len(buf)))
	copy(out[_headerLen:], buf)
	return out, nil
}

func (p ProtobufCodec) Decode(c gnet.Conn) ([]byte, error) {
	buf, _ := c.Peek(_headerLen)
	if len(buf) < _headerLen {
		return nil, ErrIncompletePacket
	}

	payloadLen := binary.BigEndian.Uint32(buf)
	msgLen := _headerLen + int(payloadLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}

	buf, _ = c.Next(msgLen)
	// todo check protobuf magic number, discard buf on error
	return buf[_headerLen:], nil
}

func (p ProtobufCodec) Unpack(buf []byte) ([]byte, error) {
	payloadLen := binary.BigEndian.Uint32(buf[:_headerLen])
	return buf[_headerLen : _headerLen+payloadLen], nil
}
