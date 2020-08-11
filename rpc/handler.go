package rpc

import (
	"fmt"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/model"
	"strconv"
	"strings"
	"time"
)

type TransferResp struct {
	Msg        string
	Total      int
	ErrInvalid int
	Latency    int64
}

func (t *TransferResp) String() string {
	s := fmt.Sprintf("TransferResp total=%d, err_invalid=%d, latency=%dms",
		t.Total, t.ErrInvalid, t.Latency)
	if t.Msg != "" {
		s = fmt.Sprintf("%s, msg=%s", s, t.Msg)
	}
	return s
}

func ReformatTag(str string) map[string]string {
	var (
		tagMap map[string]string
		tags   []string
	)

	if str == "" {
		return map[string]string{}
	}

	if strings.ContainsRune(str, ' ') {
		str = strings.Replace(str, " ", "", -1)
	}

	tagMap = make(map[string]string)

	tags = strings.Split(str, ",")
	for _, tag := range tags {
		idx := strings.IndexRune(tag, '=')
		if idx != -1 {
			tagMap[tag[:idx]] = tag[idx+1:]
		}
	}
	return tagMap
}

func RecvMetricValues(args []*model.MetricValue, reply *model.TransferResponse, source string) error {
	var (
		start time.Time
		items []*model.MetaData
	)

	start = time.Now()
	reply.Invalid = 0
	items = []*model.MetaData{}

	for _, val := range args {
		var (
			now       int64
			flag      bool
			newMetric *model.MetaData
			newValue  float64
			err       error
		)

		// 对上报数据进行校验
		if val == nil {
			reply.Invalid += 1
			continue
		}

		if val.Metric == "" || val.Endpoint == "" {
			reply.Invalid += 1
			continue
		}

		if val.Type != g.COUNTER && val.Type != g.GAUGE && val.Type != g.DERIVE {
			reply.Invalid += 1
			continue
		}

		if val.Value == "" {
			reply.Invalid += 1
			continue
		}

		if val.Step <= 0 {
			reply.Invalid += 1
			continue
		}

		if len(val.Metric)+len(val.Tags) > 510 {
			reply.Invalid += 1
			continue
		}

		// 上报数据时间处理
		now = start.Unix()
		if val.Timestamp <= 0 || val.Timestamp > now*2 {
			val.Timestamp = now
		}

		// 构建新的metric
		newMetric = &model.MetaData{
			Metric:      val.Metric,
			Endpoint:    val.Endpoint,
			Timestamp:   val.Timestamp,
			Step:        val.Step,
			CounterType: val.Type,
			Tags:        ReformatTag(val.Tags),
		}

		// metric value 类型转换
		flag = true
		switch rawData := val.Value.(type) {
		case string:
			newValue, err = strconv.ParseFloat(rawData, 64)
			if err != nil {
				flag = false
			}
		case float64:
			newValue = rawData
		case int64:
			newValue = float64(rawData)
		default:
			flag = false
		}

		if !flag {
			reply.Invalid += 1
			continue
		}

		newMetric.Value = newValue
		items = append(items, newMetric)
	}

	//cfg := g.Config()
	//
	//if cfg.Graph.Enabled {
	//	sender.Push2GraphSendQueue(items)
	//}
	//
	//if cfg.Judge.Enabled {
	//	sender.Push2JudgeSendQueue(items)
	//}
	//
	//if cfg.Tsdb.Enabled {
	//	sender.Push2TsdbSendQueue(items)
	//}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
