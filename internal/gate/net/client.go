package net

import (
	"errors"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"ppim/internal/gate/rpc"
	"sync"
	"time"
)

type (
	Client struct {
		Conn              gnet.Conn
		mtx               sync.RWMutex
		UID               string    // 用户ID
		DID               string    // 设备ID
		HeartbeatLastTime time.Time // 最后心跳时间
	}

	ClientManager struct {
		rwMtx       sync.RWMutex
		userConnMap map[string][]*Client
		connMap     map[int]*Client
	}
)

var (
	ErrConnContextNil = errors.New("客户端连接上下文为空")
)

func (c *Client) getConnContext() (*EventConnContext, error) {
	ctx, ok := c.Conn.Context().(*EventConnContext)
	if !ok {
		return nil, ErrConnContextNil
	}
	return ctx, nil
}

func (c *Client) Write(data []byte) (int, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	ctx, err := c.getConnContext()
	if err != nil {
		return 0, err
	}

	if ctx.Network == _ws {
		err = wsutil.WriteServerBinary(c.Conn, data)
		return len(data), err
	}

	buf, err := ctx.Codec.Encode(data)
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(buf)
}

func (c *Client) SetAuthResult(authed bool) error {
	ctx, err := c.getConnContext()
	if err != nil {
		return err
	}
	ctx.Authed = authed
	return nil
}

func newClientManager() *ClientManager {
	return &ClientManager{
		userConnMap: make(map[string][]*Client),
		connMap:     make(map[int]*Client),
	}
}

func (cm *ClientManager) Add(client *Client) {
	cm.rwMtx.Lock()
	defer cm.rwMtx.Unlock()

	cm.connMap[client.Conn.Fd()] = client

	ucSlice := cm.userConnMap[client.UID]
	if ucSlice == nil {
		ucSlice = make([]*Client, 0)
	}
	cm.userConnMap[client.UID] = append(ucSlice, client)

	// 登记online
	_ = rpc.Caller().Register(rpc.Context(), client.UID, client.DID)
}

func (cm *ClientManager) RemoveWithFD(fd int) {
	cm.rwMtx.Lock()
	defer cm.rwMtx.Unlock()

	client := cm.connMap[fd]
	if client == nil {
		return
	}

	// 关闭连接
	_ = client.Conn.Close()

	delete(cm.connMap, fd)

	if ucSlice := cm.userConnMap[client.UID]; len(ucSlice) > 0 {
		for i, c := range ucSlice {
			if c == client {
				cm.userConnMap[client.UID] = append(ucSlice[:i], ucSlice[i+1:]...)
			}
		}
	}

	// 注销online
	_ = rpc.Caller().UnRegister(rpc.Context(), client.UID, client.DID)
}

func (cm *ClientManager) GetWithUID(uid string) []*Client {
	cm.rwMtx.RLock()
	defer cm.rwMtx.RUnlock()

	return cm.userConnMap[uid]
}

func (cm *ClientManager) GetWithUIDAndDID(uid, did string) *Client {
	cm.rwMtx.RLock()
	defer cm.rwMtx.RUnlock()

	ucSlice := cm.userConnMap[uid]
	if len(ucSlice) > 0 {
		for i, c := range ucSlice {
			if c.DID == did {
				return ucSlice[i]
			}
		}
	}
	return nil
}

func (cm *ClientManager) GetWithFD(fd int) *Client {
	cm.rwMtx.RLock()
	defer cm.rwMtx.RUnlock()

	return cm.connMap[fd]
}
