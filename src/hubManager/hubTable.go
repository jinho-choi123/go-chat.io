package hubManager

import (
	"errors"
	"log"
	"sync"

	"go-chat.io/src/redisManager"
)

type HubTable struct {
	Mu    sync.Mutex
	Table map[uint64]*Hub
}

func NewHubTable() *HubTable {
	return &HubTable{
		Table: make(map[uint64]*Hub),
	}
}

func (ht *HubTable) Lookup(hubID uint64) (*Hub, error) {
	ht.Mu.Lock()
	hub, ok := ht.Table[hubID]
	ht.Mu.Unlock()
	if !ok {
		return nil, errors.New("cannot find the hub for hubId")
	}

	return hub, nil
}

func (ht *HubTable) Add(hubID uint64, rdm *redisManager.RedisManager) error {
	ht.Mu.Lock()

	_, ok := ht.Table[hubID]
	if ok {
		ht.Mu.Unlock()
		return errors.New("Hub already exists")
	}

	newhub := NewHub(hubID)

	go newhub.Run()
	go newhub.RedisPub(rdm)
	go newhub.RedisSub(rdm)
	log.Printf("newhub %v", newhub)

	ht.Table[hubID] = newhub

	ht.Mu.Unlock()
	return nil
}
