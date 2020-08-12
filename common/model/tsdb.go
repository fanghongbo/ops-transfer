package model

import (
	"fmt"
	"strings"
)

type TsDBItem struct {
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func (u *TsDBItem) String() string {
	return fmt.Sprintf(
		"<Metric:%s, Tags:%v, Value:%v, TS:%d>",
		u.Metric,
		u.Tags,
		u.Value,
		u.Timestamp,
	)
}

func (u *TsDBItem) TsDBString() (s string) {
	s = fmt.Sprintf("put %s %d %.3f ", u.Metric, u.Timestamp, u.Value)

	for k, v := range u.Tags {
		key := strings.ToLower(strings.Replace(k, " ", "_", -1))
		value := strings.Replace(v, " ", "_", -1)
		s += key + "=" + value + " "
	}

	return s
}
