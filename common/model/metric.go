package model

import "fmt"

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (u *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		u.Endpoint,
		u.Metric,
		u.Type,
		u.Tags,
		u.Step,
		u.Timestamp,
		u.Value,
	)
}

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

func (u *MetaData) String() string {
	return fmt.Sprintf(
		"<MetaData Endpoint:%s, Metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v>",
		u.Endpoint,
		u.Metric,
		u.Timestamp,
		u.Step,
		u.Value,
		u.Tags,
	)
}
