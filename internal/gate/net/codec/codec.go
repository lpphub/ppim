package codec

import (
	"errors"
	"github.com/panjf2000/gnet/v2"
)

type Codec[T any] interface {
	// Encode the given bytes into a codec.
	Encode(buf []byte) ([]byte, error)
	// Decode the codec into bytes.
	Decode(conn gnet.Conn) ([]T, error)
}

var (
	ErrIncompletePacket = errors.New("incomplete packet")
)
