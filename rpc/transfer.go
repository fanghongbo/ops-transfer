package rpc

import (
	"github.com/fanghongbo/ops-transfer/common/model"
)

type Transfer int

func (u *Transfer) Ping(req model.NullRpcRequest, resp *model.SimpleRpcResponse) error {
	return nil
}

func (u *Transfer) Update(args []*model.MetricValue, reply *model.TransferResponse) error {
	return RecvMetricValues(args, reply, "rpc")
}
