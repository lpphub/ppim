package util

import (
	"github.com/cespare/xxhash/v2"
	"hash/crc32"
)

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

func XXHash64(str string) uint64 {
	return xxhash.Sum64([]byte(str))
}
