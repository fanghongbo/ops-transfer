package pool

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type RpcConnPools struct {
	sync.RWMutex
	M           map[string]*ConnPool
	MaxConn     int
	MaxIdle     int
	ConnTimeout int
	CallTimeout int
}

func CreateRpcConnPools(maxConn, maxIdle, connTimeout, callTimeout int, cluster []string) *RpcConnPools {
	var (
		cp *RpcConnPools
		ct time.Duration
	)

	cp = &RpcConnPools{
		M:           make(map[string]*ConnPool),
		MaxConn:     maxConn,
		MaxIdle:     maxIdle,
		ConnTimeout: connTimeout,
		CallTimeout: callTimeout,
	}

	ct = time.Duration(cp.ConnTimeout) * time.Millisecond

	for _, address := range cluster {
		if _, exist := cp.M[address]; exist {
			continue
		}
		cp.M[address] = createOneRpcPool(address, address, ct, maxConn, maxIdle)
	}

	return cp
}

func (u *RpcConnPools) Call(addr, method string, args interface{}, resp interface{}) error {
	var (
		connPool    *ConnPool
		exist       bool
		conn        NConn
		rpcClient   *RpcClient
		callTimeout time.Duration
		done        chan error
		err         error
	)

	connPool, exist = u.Get(addr)
	if !exist {
		return fmt.Errorf("%s has no connection pool", addr)
	}

	conn, err = connPool.Fetch()
	if err != nil {
		return fmt.Errorf("%s get connection fail: conn %v, err %v. proc: %s", addr, conn, err, connPool.Proc())
	}

	rpcClient = conn.(*RpcClient)
	callTimeout = time.Duration(u.CallTimeout) * time.Millisecond

	done = make(chan error, 1)
	go func() {
		done <- rpcClient.Call(method, args, resp)
	}()

	select {
	case <-time.After(callTimeout):
		connPool.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", addr)
	case err = <-done:
		if err != nil {
			connPool.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", addr, err, connPool.Proc())
		} else {
			connPool.Release(conn)
		}
		return err
	}
}

func (u *RpcConnPools) Get(address string) (*ConnPool, bool) {
	var (
		p     *ConnPool
		exist bool
	)

	u.RLock()
	defer u.RUnlock()
	p, exist = u.M[address]
	return p, exist
}

func (u *RpcConnPools) Destroy() {
	var (
		addresses []string
	)

	u.Lock()
	defer u.Unlock()
	addresses = make([]string, 0, len(u.M))
	for address := range u.M {
		addresses = append(addresses, address)
	}

	for _, address := range addresses {
		u.M[address].Destroy()
		delete(u.M, address)
	}
}

func (u *RpcConnPools) Proc() []string {
	var proc []string

	proc = []string{}
	for _, cp := range u.M {
		proc = append(proc, cp.Proc())
	}
	return proc
}

func createOneRpcPool(name string, address string, connTimeout time.Duration, maxConn int, maxIdle int) *ConnPool {
	var p *ConnPool

	p = NewConnPool(name, address, int32(maxConn), int32(maxIdle))
	p.New = func(connName string) (NConn, error) {
		var (
			err  error
			conn net.Conn
		)

		_, err = net.ResolveTCPAddr("tcp", p.Address)
		if err != nil {
			return nil, err
		}

		conn, err = net.DialTimeout("tcp", p.Address, connTimeout)
		if err != nil {
			return nil, err
		}

		return NewRpcClient(rpc.NewClient(conn), connName), nil
	}

	return p
}

type TsDBConnPools struct {
	p           *ConnPool
	maxConn     int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

func NewTsDBConnPool(address string, maxConn, maxIdle, connTimeout, callTimeout int) *TsDBConnPools {
	return &TsDBConnPools{
		p:           newTsDBConnPool(address, maxConn, maxIdle, connTimeout),
		maxConn:     maxConn,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

func (t *TsDBConnPools) Send(data []byte) error {
	var (
		conn NConn
		cli  *struct{ net.Conn }
		done chan error
		err  error
	)

	conn, err = t.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, t.p.Proc())
	}

	cli = conn.(TsDBClient).cli

	done = make(chan error, 1)
	go func() {
		_, err = cli.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
		t.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", t.address)
	case err = <-done:
		if err != nil {
			t.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", t.address, err, t.p.Proc())
		} else {
			t.p.Release(conn)
		}
		return err
	}
}

func (t *TsDBConnPools) Destroy() {
	if t.p != nil {
		t.p.Destroy()
	}
}
