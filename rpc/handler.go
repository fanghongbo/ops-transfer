package rpc

import (
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/model"
	"github.com/fanghongbo/ops-transfer/common/proc"
	"github.com/fanghongbo/ops-transfer/sender"
	"github.com/fanghongbo/ops-transfer/utils"
	"strconv"
	"time"
)

func RecvMetricValues(args []*model.MetricValue, reply *model.TransferResponse, source string) error {
	var (
		start time.Time
		items []*model.MetaData
		cfg   *g.GlobalConfig
		cnt   int64
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
			Tags:        utils.ReformatTag(val.Tags),
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

	// statistics
	cnt = int64(len(items))
	proc.RecvCnt.IncrBy(cnt)

	if source == "rpc" {
		proc.RpcRecvCnt.IncrBy(cnt)
	}

	if source == "http" {
		proc.HttpRecvCnt.IncrBy(cnt)
	}

	cfg = g.Conf()

	// 推送到对应的队列中
	if cfg.Graph != nil && cfg.Graph.Enabled {
		sender.Push2GraphSendQueue(items)
	}

	if cfg.Judge != nil && cfg.Judge.Enabled {
		sender.Push2JudgeSendQueue(items)
	}

	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		sender.Push2TsDBSendQueue(items)
	}

	reply.Message = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
