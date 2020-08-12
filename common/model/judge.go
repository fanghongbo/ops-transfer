package model

import (
	"fmt"
	"github.com/fanghongbo/ops-transfer/utils"
)

type JudgeItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	JudgeType string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}

func (u *JudgeItem) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Value:%f, Timestamp:%d, JudgeType:%s Tags:%v>",
		u.Endpoint,
		u.Metric,
		u.Value,
		u.Timestamp,
		u.JudgeType,
		u.Tags,
	)
}

func (u *JudgeItem) PrimaryKey() string {
	return utils.Md5(utils.GetMetricPrimaryKey(u.Endpoint, u.Metric, u.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
