package chatlib

import (
	"errors"
	"fmt"
	"hash/adler32"
)

const (
	ConvSingle = "single"
	ConvGroup  = "group"
)

// GenConversationID 生成会话ID: 单聊：single|maxID@minID 群聊：group|groupID
func GenConversationID(from, to, conversationType string) (string, error) {
	switch conversationType {
	case ConvSingle:
		fromID := DigitizeUID(from)
		toID := DigitizeUID(to)
		if fromID > toID {
			return fmt.Sprintf("%s|%d@%d", conversationType, fromID, toID), nil
		}
		return fmt.Sprintf("%s|%d@%d", conversationType, toID, fromID), nil
	case ConvGroup:
		return fmt.Sprintf("%s|%s", conversationType, to), nil
	}
	return "", errors.New("unknown conversation type")
}

func DigitizeUID(uid string) uint32 {
	return adler32.Checksum([]byte(uid))
}
