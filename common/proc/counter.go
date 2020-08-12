package proc

import (
	"github.com/fanghongbo/ops-transfer/utils"
	"sync"
	"time"
)

const (
	DefaultOtherMaxSize      = 100
	DefaultSCounterQpsPeriod = 1
)

// basic counter
type SCounterBase struct {
	sync.RWMutex
	Name  string
	Cnt   int64
	Time  string
	ts    int64
	Other map[string]interface{}
}

func NewSCounterBase(name string) *SCounterBase {
	uts := time.Now().Unix()
	return &SCounterBase{Name: name, Cnt: 0, Time: utils.UnixTsFormat(uts),
		ts: uts, Other: make(map[string]interface{})}
}

func (u *SCounterBase) Get() *SCounterBase {
	u.RLock()
	defer u.RUnlock()

	return &SCounterBase{
		Name:  u.Name,
		Cnt:   u.Cnt,
		Time:  u.Time,
		ts:    u.ts,
		Other: deepCopyMap(u.Other),
	}
}

func (u *SCounterBase) SetCnt(cnt int64) {
	u.Lock()
	u.Cnt = cnt
	u.ts = time.Now().Unix()
	u.Time = utils.UnixTsFormat(u.ts)
	u.Unlock()
}

func (u *SCounterBase) Incr() {
	u.IncrBy(int64(1))
}

func (u *SCounterBase) IncrBy(incr int64) {
	u.Lock()
	u.Cnt += incr
	u.Unlock()
}

func (u *SCounterBase) PutOther(key string, value interface{}) bool {
	var (
		ret   bool
		exist bool
	)

	u.Lock()
	defer u.Unlock()

	ret = false
	_, exist = u.Other[key]
	if exist {
		u.Other[key] = value
		ret = true
	} else {
		if len(u.Other) < DefaultOtherMaxSize {
			u.Other[key] = value
			ret = true
		}
	}

	return ret
}

// counter with qps
type SCounterQps struct {
	sync.RWMutex
	Name    string
	Cnt     int64
	Qps     int64
	Time    string
	ts      int64
	lastTs  int64
	lastCnt int64
	Other   map[string]interface{}
}

func NewSCounterQps(name string) *SCounterQps {
	uts := time.Now().Unix()
	return &SCounterQps{Name: name, Cnt: 0, Time: utils.UnixTsFormat(uts), ts: uts,
		Qps: 0, lastCnt: 0, lastTs: uts, Other: make(map[string]interface{})}
}

func (u *SCounterQps) Get() *SCounterQps {
	u.Lock()
	defer u.Unlock()

	u.ts = time.Now().Unix()
	u.Time = utils.UnixTsFormat(u.ts)
	// get smooth qps value
	if u.ts-u.lastTs > DefaultSCounterQpsPeriod {
		u.Qps = int64((u.Cnt - u.lastCnt) / (u.ts - u.lastTs))
		u.lastTs = u.ts
		u.lastCnt = u.Cnt
	}

	return &SCounterQps{
		Name:    u.Name,
		Cnt:     u.Cnt,
		Qps:     u.Qps,
		Time:    u.Time,
		ts:      u.ts,
		lastTs:  u.lastTs,
		lastCnt: u.lastCnt,
		Other:   deepCopyMap(u.Other),
	}
}

func (u *SCounterQps) Incr() {
	u.IncrBy(int64(1))
}

func (u *SCounterQps) IncrBy(incr int64) {
	u.Lock()
	u.incrBy(incr)
	u.Unlock()
}

func (u *SCounterQps) PutOther(key string, value interface{}) bool {
	var (
		ret   bool
		exist bool
	)

	u.Lock()
	defer u.Unlock()

	ret = false
	_, exist = u.Other[key]
	if exist {
		u.Other[key] = value
		ret = true
	} else {
		if len(u.Other) < DefaultOtherMaxSize {
			u.Other[key] = value
			ret = true
		}
	}

	return ret
}

func (u *SCounterQps) incrBy(incr int64) {
	u.Cnt += incr
}

func deepCopyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for key, val := range src {
		dst[key] = val
	}
	return dst
}
