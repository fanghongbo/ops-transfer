package model

import "fmt"

// code == 0 => success
// code == 1 => bad request
type SimpleRpcResponse struct {
	Code int `json:"code"`
}

func (u *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", u.Code)
}

type NullRpcRequest struct {
}
