package net

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type Client struct {
	Conn   *gnet.Conn
	Authed bool   // 是否认证
	UID    string // 用户ID
	DID    string // 设备ID
}

type ClientManager struct {
	rwMtx         sync.RWMutex
	userClientMap map[string][]*Client
}

func (cm *ClientManager) Add(client *Client) {
	cm.rwMtx.Lock()
	defer cm.rwMtx.Unlock()

	if ucSlice, ok := cm.userClientMap[client.UID]; !ok {
		cm.userClientMap[client.UID] = []*Client{client}
	} else {
		ucSlice = append(ucSlice, client)
	}
}

func (cm *ClientManager) Remove(client *Client) {
	cm.rwMtx.Lock()
	defer cm.rwMtx.Unlock()

	if ucSlice, ok := cm.userClientMap[client.UID]; ok && len(ucSlice) > 0 {
		for i, c := range ucSlice {
			if c == client {
				ucSlice = append(ucSlice[:i], ucSlice[i+1:]...)
				cm.userClientMap[client.UID] = ucSlice
			}
		}
	}
}

func (cm *ClientManager) GetWithUID(uid string) []*Client {
	cm.rwMtx.RLock()
	defer cm.rwMtx.RUnlock()

	return cm.userClientMap[uid]
}

func (cm *ClientManager) GetWithUIDAndDID(uid, did string) *Client {
	cm.rwMtx.RLock()
	defer cm.rwMtx.RUnlock()

	if ucSlice, ok := cm.userClientMap[uid]; ok && len(ucSlice) > 0 {
		for i, c := range ucSlice {
			if c.DID == did {
				return ucSlice[i]
			}
		}
	}
	return nil
}
