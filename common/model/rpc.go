package model

import "fmt"

type RpcResponse struct {
	Code int `json:"code"`  // 0 success; 1 fail
}

func (u *RpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", u.Code)
}

type NullRpcRequest struct {
}
