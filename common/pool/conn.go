package pool

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"io"
	"sync"
	"time"
)

var ErrMaxConn = fmt.Errorf("maximum connections reached")

type NConn interface {
	io.Closer
	Name() string
	Closed() bool
}

// conn pool
type ConnPool struct {
	sync.RWMutex

	Name    string
	Address string
	MaxConn int32
	MaxIdle int32
	Cnt     int64

	New func(name string) (NConn, error)

	active int32
	free   []NConn
	all    map[string]NConn
}

func NewConnPool(name string, address string, maxConn int32, maxIdle int32) *ConnPool {
	return &ConnPool{Name: name, Address: address, MaxConn: maxConn, MaxIdle: maxIdle, Cnt: 0, all: make(map[string]NConn)}
}

func (u *ConnPool) Proc() string {
	u.RLock()
	defer u.RUnlock()

	return fmt.Sprintf("Name:%s,Cnt:%d,active:%d,all:%d,free:%d",
		u.Name, u.Cnt, u.active, len(u.all), len(u.free))
}

func (u *ConnPool) Fetch() (NConn, error) {
	u.Lock()
	defer u.Unlock()

	// get from free
	conn := u.fetchFree()
	if conn != nil {
		return conn, nil
	}

	if u.overMax() {
		return nil, ErrMaxConn
	}

	// create new pool
	conn, err := u.newConn()
	if err != nil {
		return nil, err
	}

	u.incrActive()
	return conn, nil
}

func (u *ConnPool) Release(conn NConn) {
	u.Lock()
	defer u.Unlock()

	if u.overMaxIdle() {
		u.deleteConn(conn)
		u.decrActive()
	} else {
		u.addFree(conn)
	}
}

func (u *ConnPool) ForceClose(conn NConn) {
	u.Lock()
	defer u.Unlock()

	u.deleteConn(conn)
	u.decrActive()
}

func (u *ConnPool) Destroy() {
	u.Lock()
	defer u.Unlock()

	for _, conn := range u.free {
		if conn != nil && !conn.Closed() {
			if err := conn.Close(); err != nil {
				dlog.Error(err)
			}
		}
	}

	for _, conn := range u.all {
		if conn != nil && !conn.Closed() {
			if err := conn.Close(); err != nil {
				dlog.Error(err)
			}
		}
	}

	u.active = 0
	u.free = []NConn{}
	u.all = map[string]NConn{}
}

// internal, concurrently unsafe
func (u *ConnPool) newConn() (NConn, error) {
	name := fmt.Sprintf("%s_%d_%d", u.Name, u.Cnt, time.Now().Unix())
	conn, err := u.New(name)
	if err != nil {
		if conn != nil {
			if err := conn.Close(); err != nil {
				dlog.Error(err)
			}
		}
		return nil, err
	}

	u.Cnt++
	u.all[conn.Name()] = conn
	return conn, nil
}

func (u *ConnPool) deleteConn(conn NConn) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			dlog.Error(err)
		}
		delete(u.all, conn.Name())
	}

}

func (u *ConnPool) addFree(conn NConn) {
	u.free = append(u.free, conn)
}

func (u *ConnPool) fetchFree() NConn {
	if len(u.free) == 0 {
		return nil
	}

	conn := u.free[0]
	u.free = u.free[1:]
	return conn
}

func (u *ConnPool) incrActive() {
	u.active += 1
}

func (u *ConnPool) decrActive() {
	u.active -= 1
}

func (u *ConnPool) overMax() bool {
	return u.active >= u.MaxConn
}

func (u *ConnPool) overMaxIdle() bool {
	return int32(len(u.free)) >= u.MaxIdle
}
