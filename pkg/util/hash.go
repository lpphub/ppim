package util

import (
	"github.com/cespare/xxhash/v2"
	"github.com/spaolacci/murmur3"
	"hash/crc32"
)

func Murmur32(str string) uint32 {
	return murmur3.Sum32([]byte(str))
}

func Murmur64(str string) uint64 {
	return murmur3.Sum64([]byte(str))
}

func XXHash64(str string) uint64 {
	return xxhash.Sum64([]byte(str))
}

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
