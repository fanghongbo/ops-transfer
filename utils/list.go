package utils

import (
	"container/list"
	"sync"
)

type List struct {
	sync.RWMutex
	L *list.List
}

func NewSafeList() *List {
	return &List{L: list.New()}
}

func (u *List) PushFront(v interface{}) *list.Element {
	u.Lock()
	defer u.Unlock()

	return u.L.PushFront(v)
}

func (u *List) PushFrontBatch(vs []interface{}) {
	u.Lock()
	defer u.Unlock()

	for _, item := range vs {
		u.L.PushFront(item)
	}
}

func (u *List) PopBack() interface{} {
	u.Lock()

	if elem := u.L.Back(); elem != nil {
		item := u.L.Remove(elem)
		u.Unlock()
		return item
	}
	u.Unlock()
	return nil
}

func (u *List) PopBackBy(max int) []interface{} {
	var (
		count int
		items []interface{}
	)

	u.Lock()

	count = u.len()
	if count == 0 {
		u.Unlock()
		return []interface{}{}
	}

	if count > max {
		count = max
	}

	items = make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := u.L.Remove(u.L.Back())
		items = append(items, item)
	}

	u.Unlock()
	return items
}

func (u *List) PopBackAll() []interface{} {
	var (
		count int
		items []interface{}
	)

	u.Lock()

	count = u.len()
	if count == 0 {
		u.Unlock()
		return []interface{}{}
	}

	items = make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		item := u.L.Remove(u.L.Back())
		items = append(items, item)
	}

	u.Unlock()
	return items
}

func (u *List) Remove(e *list.Element) interface{} {
	u.Lock()
	defer u.Unlock()
	return u.L.Remove(e)
}

func (u *List) RemoveAll() {
	u.Lock()
	u.L = list.New()
	u.Unlock()
}

func (u *List) FrontAll() []interface{} {
	var (
		count int
		items []interface{}
	)

	u.RLock()
	defer u.RUnlock()

	count = u.len()
	if count == 0 {
		return []interface{}{}
	}

	items = make([]interface{}, 0, count)
	for e := u.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}

func (u *List) BackAll() []interface{} {
	u.RLock()
	defer u.RUnlock()

	count := u.len()
	if count == 0 {
		return []interface{}{}
	}

	items := make([]interface{}, 0, count)
	for e := u.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
	return items
}

func (u *List) Front() interface{} {
	u.RLock()

	if f := u.L.Front(); f != nil {
		u.RUnlock()
		return f.Value
	}

	u.RUnlock()
	return nil
}

func (u *List) Len() int {
	u.RLock()
	defer u.RUnlock()
	return u.len()
}

func (u *List) len() int {
	return u.L.Len()
}

// SafeList with Limited Size
type ListLimited struct {
	maxSize int
	SL      *List
}

func NewSafeListLimited(maxSize int) *ListLimited {
	return &ListLimited{SL: NewSafeList(), maxSize: maxSize}
}

func (u *ListLimited) PopBack() interface{} {
	return u.SL.PopBack()
}

func (u *ListLimited) PopBackBy(max int) []interface{} {
	return u.SL.PopBackBy(max)
}

func (u *ListLimited) PushFront(v interface{}) bool {
	if u.SL.Len() >= u.maxSize {
		return false
	}

	u.SL.PushFront(v)
	return true
}

func (u *ListLimited) PushFrontBatch(vs []interface{}) bool {
	if u.SL.Len() >= u.maxSize {
		return false
	}

	u.SL.PushFrontBatch(vs)
	return true
}

func (u *ListLimited) PushFrontViolently(v interface{}) bool {
	u.SL.PushFront(v)
	if u.SL.Len() > u.maxSize {
		u.SL.PopBack()
	}

	return true
}

func (u *ListLimited) RemoveAll() {
	u.SL.RemoveAll()
}

func (u *ListLimited) Front() interface{} {
	return u.SL.Front()
}

func (u *ListLimited) FrontAll() []interface{} {
	return u.SL.FrontAll()
}

func (u *ListLimited) Len() int {
	return u.SL.Len()
}
