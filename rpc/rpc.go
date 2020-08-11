package rpc

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

func Start() {
	var (
		addr     string
		tcpAddr  *net.TCPAddr
		server   *rpc.Server
		listener *net.TCPListener
		err      error
	)

	if g.Conf().Rpc == nil || !g.Conf().Rpc.Enabled {
		dlog.Warning("rpc is disable")
		return
	}

	addr = g.Conf().Rpc.Listen

	tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		dlog.Fatalf("resolve tcp addr err: %s", err)
	}

	server = rpc.NewServer()
	if err = server.Register(new(Transfer)); err != nil {
		dlog.Fatalf("register hbs err: %s", err)
	}

	listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		dlog.Fatal(err)
	} else {
		dlog.Infof("listening %s", addr)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			dlog.Errorf("listener accept err:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}

		if err = conn.SetKeepAlive(true); err != nil {
			dlog.Errorf("set rpc keep alive err: %s", err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
