package codec

import (
	"encoding/binary"
	"github.com/panjf2000/gnet/v2"
)

// Frame
// * BEFORE (300 bytes)              AFTER (304 bytes)
// * +---------------+               +--------+---------------+
// * | Protobuf Data |-------------->| Length | Protobuf Data |
// * | (300 bytes)   |               | 0x12C  | (300 bytes)   |
// * +---------------+               +--------+---------------+
type Frame struct {
	Body []byte
}

type FrameCodec struct {
}

const (
	_header = 4
)

func (f *FrameCodec) Encode(buf []byte) ([]byte, error) {
	out := make([]byte, _header+len(buf))
	binary.BigEndian.PutUint32(out, uint32(len(buf)))
	copy(out[_header:], buf)
	return out, nil
}

func (f *FrameCodec) Decode(c gnet.Conn) (frames []Frame, err error) {
	for {
		if c.InboundBuffered() < _header {
			break
		}

		// 读取消息头
		headerBuf, _ := c.Peek(_header)
		bodyLen := binary.BigEndian.Uint32(headerBuf)
		msgLen := _header + int(bodyLen)

		// 检查是否有足够的字节读取完整消息
		if c.InboundBuffered() < msgLen {
			break
		}

		// 读取完整消息
		msgBuf, _ := c.Next(msgLen)

		// 提取消息体并添加到结果中
		frames = append(frames, Frame{
			Body: msgBuf[_header:],
		})
	}

	if len(frames) == 0 {
		return frames, ErrIncompletePacket
	}
	return
}
