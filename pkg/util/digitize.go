package util

import (
	"hash/adler32"
	"hash/crc32"
)

func DigitizeWithAdler32(str string) uint32 {
	return adler32.Checksum([]byte(str))
}

func DigitizeWithCRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
