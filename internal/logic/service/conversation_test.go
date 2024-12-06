package service

import (
	"hash/adler32"
	"hash/crc32"
	"testing"
)

func TestConversationSrv_IndexConversation(t *testing.T) {
	data := []byte("1235567886")
	val := adler32.Checksum(data)
	t.Log(val)

	val2 := crc32.ChecksumIEEE(data)
	t.Log(val2)

}
