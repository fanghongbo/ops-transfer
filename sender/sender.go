package sender

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/model"
	"github.com/fanghongbo/ops-transfer/common/proc"
	"github.com/fanghongbo/ops-transfer/utils"
)

var (
	MinStep int // 最小上报周期,单位 sec
)

func Start() {
	MinStep = g.Conf().MinStep

	initConnPools()
	initSendQueues()
	initNodeRings()
	startSendTasks()
	startSenderCron()
}

// 将数据 打入 某个Judge的发送缓存队列, 具体是哪一个Judge 由一致性哈希 决定
func Push2JudgeSendQueue(items []*model.MetaData) {
	for _, item := range items {
		var (
			pk        string
			node      string
			err       error
			step      int
			ts        int64
			judgeItem *model.JudgeItem
			queue     *utils.ListLimited
			isSuccess bool
		)

		pk = item.PK()
		node, err = JudgeNodeRing.GetNode(pk)
		if err != nil {
			dlog.Errorf("get hash node err: %s", err)
			continue
		}

		// align ts
		step = int(item.Step)
		if step < MinStep {
			step = MinStep
		}

		ts = alignTs(item.Timestamp, int64(step))

		judgeItem = &model.JudgeItem{
			Endpoint:  item.Endpoint,
			Metric:    item.Metric,
			Value:     item.Value,
			Timestamp: ts,
			JudgeType: item.CounterType,
			Tags:      item.Tags,
		}

		queue = JudgeQueues[node]
		isSuccess = queue.PushFront(judgeItem)

		// statistics
		if !isSuccess {
			proc.SendToJudgeDropCnt.Incr()
		}
	}
}

// 将数据 打入 某个Graph的发送缓存队列, 具体是哪一个Graph 由一致性哈希 决定
func Push2GraphSendQueue(items []*model.MetaData) {
	for _, item := range items {
		var (
			graphItem   *model.GraphItem
			pk          string
			node        string
			clusterNode *g.ClusterNode
			errCnt      int
			err         error
		)

		graphItem, err = convert2GraphItem(item)
		if err != nil {
			dlog.Errorf("convert to graph item err: %s", err)
			continue
		}

		pk = item.PK()

		// statistics. 为了效率,放到了这里,因此只有graph是enbale时才能trace
		proc.RecvDataTrace.Trace(pk, item)
		proc.RecvDataFilter.Filter(pk, item.Value, item)

		node, err = GraphNodeRing.GetNode(pk)
		if err != nil {
			dlog.Errorf("get hash node err: %s", err)
			continue
		}

		clusterNode = g.GraphCluster[node]
		errCnt = 0
		for _, addr := range clusterNode.Addrs {
			Q := GraphQueues[node+addr]
			if !Q.PushFront(graphItem) {
				errCnt += 1
			}
		}

		// statistics
		if errCnt > 0 {
			proc.SendToGraphDropCnt.Incr()
		}
	}
}

// 打到Graph的数据,要根据rrd tool的设定 来限制 step、counterType、timestamp
func convert2GraphItem(d *model.MetaData) (*model.GraphItem, error) {
	var item *model.GraphItem

	item = &model.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < MinStep {
		item.Step = MinStep
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == g.GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == g.COUNTER {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == g.DERIVE {
		item.DsType = g.DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not supported counter type")
	}

	item.Timestamp = alignTs(item.Timestamp, int64(item.Step)) //item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

// 将原始数据入到ts db发送缓存队列
func Push2TsDBSendQueue(items []*model.MetaData) {
	for _, item := range items {
		var (
			newItem   *model.TsDBItem
			isSuccess bool
		)

		newItem = convert2TsDBItem(item)
		isSuccess = TsDBQueue.PushFront(newItem)

		if !isSuccess {
			proc.SendToTsDBDropCnt.Incr()
		}
	}
}

// 转化为ts db格式
func convert2TsDBItem(d *model.MetaData) *model.TsDBItem {
	var t model.TsDBItem

	t = model.TsDBItem{Tags: make(map[string]string)}

	for k, v := range d.Tags {
		t.Tags[k] = v
	}

	t.Tags["endpoint"] = d.Endpoint
	t.Metric = d.Metric
	t.Timestamp = d.Timestamp
	t.Value = d.Value
	return &t
}

func alignTs(ts int64, period int64) int64 {
	return ts - ts%period
}
