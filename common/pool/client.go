package pool

import (
	"net"
	"net/rpc"
	"time"
)

type RpcClient struct {
	cli  *rpc.Client
	name string
}

func (u RpcClient) Name() string {
	return u.name
}

func (u RpcClient) Closed() bool {
	return u.cli == nil
}

func (u RpcClient) Close() error {
	if u.cli != nil {
		err := u.cli.Close()
		u.cli = nil
		return err
	}
	return nil
}

func (u RpcClient) Call(method string, args interface{}, reply interface{}) error {
	return u.cli.Call(method, args, reply)
}

func NewRpcClient(cli *rpc.Client, name string) *RpcClient {
	return &RpcClient{cli: cli, name: name}
}

func NewRpcClientWithCodec(codec rpc.ClientCodec, name string) *RpcClient {
	return &RpcClient{
		cli:  rpc.NewClientWithCodec(codec),
		name: name,
	}
}

type TsDBClient struct {
	cli  *struct{ net.Conn }
	name string
}

func (u TsDBClient) Name() string {
	return u.name
}

func (u TsDBClient) Closed() bool {
	return u.cli.Conn == nil
}

func (u TsDBClient) Close() error {
	if u.cli != nil {
		err := u.cli.Close()
		u.cli.Conn = nil
		return err
	}
	return nil
}

func newTsDBConnPool(address string, maxConn int, maxIdle int, connTimeout int) *ConnPool {
	var pool *ConnPool

	pool = NewConnPool("tsdb", address, int32(maxConn), int32(maxIdle))
	pool.New = func(name string) (NConn, error) {
		var (
			err  error
			conn net.Conn
		)

		_, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		conn, err = net.DialTimeout("tcp", address, time.Duration(connTimeout)*time.Millisecond)
		if err != nil {
			return nil, err
		}

		return TsDBClient{
			cli:  &struct{ net.Conn }{conn},
			name: name,
		}, nil
	}

	return pool
}
