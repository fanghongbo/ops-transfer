package model

import (
	"fmt"
	"github.com/fanghongbo/ops-transfer/utils"
)

// DsType 即 rrd 中的 data source 的类型：GAUGE | COUNTER | DERIVE
type GraphItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	DsType    string            `json:"dstype"`
	Step      int               `json:"step"`
	Heartbeat int               `json:"heartbeat"`
	Min       string            `json:"min"`
	Max       string            `json:"max"`
}

func (u *GraphItem) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Tags:%v, Value:%v, TS:%d %v DsType:%s, Step:%d, Heartbeat:%d, Min:%s, Max:%s>",
		u.Endpoint,
		u.Metric,
		u.Tags,
		u.Value,
		u.Timestamp,
		utils.UnixTsFormat(u.Timestamp),
		u.DsType,
		u.Step,
		u.Heartbeat,
		u.Min,
		u.Max,
	)
}

func (u *GraphItem) PrimaryKey() string {
	return utils.GetMetricPrimaryKey(u.Endpoint, u.Metric, u.Tags)
}

func (u *GraphItem) Checksum() string {
	return utils.Md5(utils.GetMetricPrimaryKey(u.Endpoint, u.Metric, u.Tags))
}

func (u *GraphItem) UUID() string {
	return utils.ReformatMetricUniqueString(u.Endpoint, u.Metric, u.Tags, u.DsType, u.Step)
}
