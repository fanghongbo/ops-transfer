package proc

import (
	"container/list"
	"fmt"
	"math"
	"sync"
)

type DataFilter struct {
	sync.RWMutex
	MaxSize   int
	Name      string
	PK        string
	Opt       string
	Threshold float64
	L         *list.List
}

func NewDataFilter(name string, maxSize int) *DataFilter {
	return &DataFilter{L: list.New(), Name: name, MaxSize: maxSize}
}

func (u *DataFilter) SetFilter(pk string, opt string, threshhold float64) error {
	u.Lock()
	defer u.Unlock()

	if !legalOpt(opt) {
		return fmt.Errorf("bad opt: %s", opt)
	}

	// rm old caches when filter's pk changed
	if u.PK != pk {
		u.L = list.New()
	}
	u.PK = pk
	u.Opt = opt
	u.Threshold = threshhold

	return nil
}

// proposed that there were few traced items
func (u *DataFilter) Filter(pk string, val float64, v interface{}) {
	u.RLock()
	if u.PK != pk {
		u.RUnlock()
		return
	}
	u.RUnlock()

	// we could almost not step here, so we get few wlock
	u.Lock()
	defer u.Unlock()
	if compute(u.Opt, val, u.Threshold) {
		u.L.PushFront(v)
		if u.L.Len() > u.MaxSize {
			u.L.Remove(u.L.Back())
		}
	}
}

func (u *DataFilter) GetAllFiltered() []interface{} {
	u.RLock()
	defer u.RUnlock()

	items := make([]interface{}, 0)
	for e := u.L.Front(); e != nil; e = e.Next() {
		items = append(items, e)
	}

	return items
}

// internal
const (
	MinPositiveFloat64 = 0.000001
	MaxNegativeFloat64 = -0.000001
)

func compute(opt string, left float64, right float64) bool {
	switch opt {
	case "eq":
		return math.Abs(left-right) < MinPositiveFloat64
	case "ne":
		return math.Abs(left-right) >= MinPositiveFloat64
	case "gt":
		return (left - right) > MinPositiveFloat64
	case "lt":
		return (left - right) < MaxNegativeFloat64
	default:
		return false
	}
}

func legalOpt(opt string) bool {
	switch opt {
	case "eq", "ne", "gt", "lt":
		return true
	default:
		return false
	}
}
