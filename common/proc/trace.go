package proc

import (
	"container/list"
	"sync"
)

type DataTrace struct {
	sync.RWMutex
	MaxSize int
	Name    string
	PK      string
	L       *list.List
}

func NewDataTrace(name string, maxSize int) *DataTrace {
	return &DataTrace{L: list.New(), Name: name, MaxSize: maxSize}
}

func (u *DataTrace) SetPK(pk string) {
	u.Lock()
	defer u.Unlock()

	// rm old caches when trace's pk changed
	if u.PK != pk {
		u.L = list.New()
	}
	u.PK = pk
}

// proposed that there were few traced items
func (u *DataTrace) Trace(pk string, v interface{}) {
	u.RLock()
	if u.PK != pk {
		u.RUnlock()
		return
	}

	u.RUnlock()

	// we could almost not step here, so we get few wlock
	u.Lock()
	defer u.Unlock()
	u.L.PushFront(v)
	if u.L.Len() > u.MaxSize {
		u.L.Remove(u.L.Back())
	}
}

func (u *DataTrace) GetAllTraced() []interface{} {
	var items []interface{}

	u.RLock()
	defer u.RUnlock()

	items = make([]interface{}, 0)
	for e := u.L.Front(); e != nil; e = e.Next() {
		items = append(items, e)
	}

	return items
}
