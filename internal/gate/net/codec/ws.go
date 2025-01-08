package codec

import (
	"bytes"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lpphub/golib/logger"
	"github.com/panjf2000/gnet/v2"
	"github.com/pkg/errors"
	"io"
)

type WsCodec struct {
	upgraded bool         // 链接是否升级
	Buf      bytes.Buffer // 从实际socket中读取到的数据缓存
	WsMsgBuf wsMessageBuf // ws 消息缓存
}

type wsMessageBuf struct {
	curHeader *ws.Header
	cachedBuf bytes.Buffer
}

type readWrite struct {
	io.Reader
	io.Writer
}

func (w *WsCodec) upgrade(c gnet.Conn) error {
	if w.upgraded {
		return nil
	}

	buf := &w.Buf
	tmpReader := bytes.NewReader(buf.Bytes())
	oldLen := tmpReader.Len()

	hs, err := ws.Upgrade(readWrite{tmpReader, c})
	skipN := oldLen - tmpReader.Len()
	if err != nil {
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) { //数据不完整，不跳过 Buf 中的 skipN 字节（此时 Buf 中存放的仅是部分 "handshake data" bytes），下次再尝试读取
			return nil
		}
		buf.Next(skipN)
		logger.Log().Error().Msgf("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		return err
	}
	buf.Next(skipN)
	logger.Log().Info().Msgf("conn[%v] upgrade websocket protocol! Handshake: %v", c.RemoteAddr().String(), hs)

	w.upgraded = true
	return nil
}
func (w *WsCodec) ReadBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size)
	read, err := c.Read(buf)
	if err != nil {
		logger.Log().Error().Msgf("read err! %v", err)
		return gnet.Close
	}
	if read < size {
		logger.Log().Error().Msgf("read bytes len err! size: %d read: %d", size, read)
		return gnet.Close
	}
	w.Buf.Write(buf)
	return gnet.None
}

func (w *WsCodec) Decode(c gnet.Conn) (outs []wsutil.Message, err error) {
	if err = w.upgrade(c); err != nil {
		return nil, err
	}
	if w.Buf.Len() <= 0 {
		return
	}

	messages, err := w.readWsMessages()
	if err != nil {
		logger.Log().Error().Msgf("Error reading message! %v", err)
		return nil, err
	}
	if messages == nil || len(messages) <= 0 { //没有读到完整数据 不处理
		return
	}
	for _, message := range messages {
		if message.OpCode.IsControl() {
			err = wsutil.HandleClientControlMessage(c, message)
			if err != nil {
				return
			}
			continue
		}
		if message.OpCode == ws.OpText || message.OpCode == ws.OpBinary {
			outs = append(outs, message)
		}
	}
	return
}

func (w *WsCodec) readWsMessages() (messages []wsutil.Message, err error) {
	msgBuf := &w.WsMsgBuf
	in := &w.Buf
	for {
		// 从 in 中读出 header，并将 header bytes 写入 msgBuf.cachedBuf
		if msgBuf.curHeader == nil {
			if in.Len() < ws.MinHeaderSize { //头长度至少是2
				return
			}
			var head ws.Header
			if in.Len() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(in)
				if err != nil {
					return messages, err
				}
			} else { //有可能不完整，构建新的 reader 读取 head，读取成功才实际对 in 进行读操作
				tmpReader := bytes.NewReader(in.Bytes())
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) { //数据不完整
						return messages, nil
					}
					in.Next(skipN)
					return
				}
				in.Next(skipN)
			}

			msgBuf.curHeader = &head
			err = ws.WriteHeader(&msgBuf.cachedBuf, head)
			if err != nil {
				return
			}
		}
		dataLen := (int)(msgBuf.curHeader.Length)
		// 从 in 中读出 data，并将 data bytes 写入 msgBuf.cachedBuf
		if dataLen > 0 {
			if in.Len() < dataLen { //数据不完整
				logger.Log().Warn().Msgf("incomplete data: dl=%d, rl=%d", dataLen, in.Len())
				return
			}

			_, err = io.CopyN(&msgBuf.cachedBuf, in, int64(dataLen))
			if err != nil {
				return
			}
		}
		if msgBuf.curHeader.Fin { //当前 header 已经是一个完整消息
			messages, err = wsutil.ReadClientMessage(&msgBuf.cachedBuf, messages)
			if err != nil {
				return
			}
			msgBuf.cachedBuf.Reset()
		} else {
			logger.Log().Error().Msgf("The data is split into multiple frames")
		}
		msgBuf.curHeader = nil
	}
}
