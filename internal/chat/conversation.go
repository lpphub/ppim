package chat

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
)

const (
	ConvSingle = "single"
	ConvGroup  = "group"
)

// GenConversationID 生成会话ID: 单聊：single|maxID@minID 群聊：group|groupID
func GenConversationID(from, to, conversationType string) (string, error) {
	switch conversationType {
	case ConvSingle:
		fromID := cast.ToInt64(from)
		toID := cast.ToInt64(to)
		if fromID > toID {
			return fmt.Sprintf("%s|%d@%d", conversationType, fromID, toID), nil
		}
		return fmt.Sprintf("%s|%d@%d", conversationType, toID, fromID), nil
	case ConvGroup:
		return fmt.Sprintf("%s|%s", conversationType, to), nil
	}
	return "", errors.New("unknown conversation type")
}
