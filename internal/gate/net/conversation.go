package net

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
)

const (
	convSingle = "single"
	convGroup  = "group"
)

// GenConversationID 生成会话ID: 单聊：single|maxID@minID 群聊：group|groupID
func GenConversationID(from, to, conversationType string) (string, error) {
	switch conversationType {
	case convSingle:
		fromID := cast.ToInt64(from)
		toID := cast.ToInt64(to)
		if fromID > toID {
			return fmt.Sprintf("%s|%d@%d", conversationType, fromID, toID), nil
		}
		return fmt.Sprintf("%s|%d@%d", conversationType, toID, fromID), nil
	case convGroup:
		return fmt.Sprintf("%s|%s", conversationType, to), nil
	}
	return "", errors.New("unknown conversation type")
}
