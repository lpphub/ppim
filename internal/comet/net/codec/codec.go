package codec

import "github.com/panjf2000/gnet/v2"

type Codec interface {
	// Encode the given bytes into a codec.
	Encode(buf []byte) ([]byte, error)
	// Decode the codec into bytes.
	Decode(conn gnet.Conn) ([]byte, error)
	// Unpack the codec into bytes.
	Unpack(buf []byte) ([]byte, error)
}
