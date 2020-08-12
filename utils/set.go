package utils

import "sync"

type StringSet struct {
	M map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{
		M: make(map[string]struct{}),
	}
}

func (u *StringSet) Add(elt string) *StringSet {
	u.M[elt] = struct{}{}
	return u
}

func (u *StringSet) Exists(elt string) bool {
	var exist bool

	_, exist = u.M[elt]
	return exist
}

func (u *StringSet) Delete(elt string) {
	delete(u.M, elt)
}

func (u *StringSet) Clear() {
	u.M = make(map[string]struct{})
}

func (u *StringSet) ToSlice() []string {
	var (
		count int
		data  []string
	)

	count = len(u.M)
	if count == 0 {
		return []string{}
	}

	data = make([]string, count)

	i := 0
	for elt := range u.M {
		data[i] = elt
		i++
	}

	return data
}

type SafeSet struct {
	sync.RWMutex
	M map[string]bool
}

func NewSafeSet() *SafeSet {
	return &SafeSet{
		M: make(map[string]bool),
	}
}

func (u *SafeSet) Add(key string) {
	u.Lock()
	defer u.Unlock()

	u.M[key] = true

}

func (u *SafeSet) Remove(key string) {
	u.Lock()
	defer u.Unlock()
	delete(u.M, key)

}

func (u *SafeSet) Clear() {
	u.Lock()
	defer u.Unlock()
	u.M = make(map[string]bool)

}

func (u *SafeSet) Contains(key string) bool {
	var exist bool

	u.RLock()
	defer u.RUnlock()

	_, exist = u.M[key]
	return exist
}

func (u *SafeSet) Size() int {
	u.RLock()
	defer u.RUnlock()
	return len(u.M)
}

func (u *SafeSet) ToSlice() []string {
	var data []string

	u.RLock()
	defer u.RUnlock()

	count := len(u.M)
	if count == 0 {
		return []string{}
	}

	data = []string{}
	for key := range u.M {
		data = append(data, key)
	}

	return data
}
