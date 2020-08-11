package http

import (
	"encoding/json"
	"github.com/fanghongbo/dlog"
	"net"
	"net/http"
	"strings"
)

func IsLocalRequest(r *http.Request) bool {
	var (
		remoteAddr []string
		localAddr  []net.Addr
		err        error
	)

	remoteAddr = strings.Split(r.RemoteAddr, ":")

	localAddr, err = net.InterfaceAddrs()
	if err != nil {
		dlog.Errorf("get local address err: %s", err)
		return false
	}

	for _, item := range localAddr {
		if ipNet, ok := item.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil {
				if remoteAddr[0] == ipNet.IP.String() {
					return true
				}
			}
		}
	}

	return false
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	var (
		bs  []byte
		err error
	)

	bs, err = json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if _, err = w.Write(bs); err != nil {
		dlog.Errorf("http response err: %s", err)
	}
}
