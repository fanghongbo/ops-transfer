package sender

import (
	"bytes"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/model"
	"github.com/fanghongbo/ops-transfer/common/proc"
	"github.com/fanghongbo/ops-transfer/utils"
	"time"
)

// send
const (
	DefaultSendTaskSleepInterval = time.Millisecond * 50 //默认睡眠间隔为50ms
)

func startSendTasks() {
	var (
		cfg             *g.GlobalConfig
		judgeConcurrent int
		graphConcurrent int
		tsDBConcurrent  int
	)

	cfg = g.Conf()

	// init semaphore
	judgeConcurrent = cfg.Judge.MaxConn
	graphConcurrent = cfg.Graph.MaxConn
	tsDBConcurrent = cfg.TsDB.MaxConn

	// init send go-routines
	for node := range cfg.Judge.Cluster {
		queue := JudgeQueues[node]
		go forward2JudgeTask(queue, node, judgeConcurrent)
	}

	for node, item := range g.GraphCluster {
		for _, addr := range item.Addrs {
			queue := GraphQueues[node+addr]
			go forward2GraphTask(queue, node, addr, graphConcurrent)
		}
	}

	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		go forward2TsDBTask(tsDBConcurrent)
	}
}

// Judge定时任务, 将 Judge发送缓存中的数据 通过rpc连接池 发送到Judge
func forward2JudgeTask(Q *utils.ListLimited, node string, concurrent int) {
	var (
		batch int
		addr  string
		lock  *utils.Semaphore
	)

	batch = g.Conf().Judge.Batch // 一次发送,最多batch条数据
	addr = g.Conf().Judge.Cluster[node]
	lock = utils.NewSemaphore(concurrent)

	for {
		var (
			items      []interface{}
			count      int
			judgeItems []*model.JudgeItem
		)

		items = Q.PopBackBy(batch)
		count = len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		judgeItems = make([]*model.JudgeItem, count)
		for i := 0; i < count; i++ {
			judgeItems[i] = items[i].(*model.JudgeItem)
		}

		//	同步Call + 有限并发 进行发送
		lock.Acquire()
		go func(addr string, judgeItems []*model.JudgeItem, count int) {
			var (
				err    error
				resp   *model.RpcResponse
				sendOk bool
			)

			defer lock.Release()

			resp = &model.RpcResponse{}
			sendOk = false

			for i := 0; i < 3; i++ { //最多重试3次
				err = JudgeConnPools.Call(addr, "Judge.Send", judgeItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				dlog.Infof("send judge %s:%s fail: %v", node, addr, err)
				proc.SendToJudgeFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToJudgeCnt.IncrBy(int64(count))
			}
		}(addr, judgeItems, count)
	}
}

// Graph定时任务, 将 Graph发送缓存中的数据 通过rpc连接池 发送到Graph
func forward2GraphTask(Q *utils.ListLimited, node string, addr string, concurrent int) {
	var (
		batch int
		lock  *utils.Semaphore
	)

	batch = g.Conf().Graph.Batch // 一次发送,最多batch条数据
	lock = utils.NewSemaphore(concurrent)

	for {
		var (
			items      []interface{}
			count      int
			graphItems []*model.GraphItem
		)

		items = Q.PopBackBy(batch)
		count = len(items)
		if count == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		graphItems = make([]*model.GraphItem, count)
		for i := 0; i < count; i++ {
			graphItems[i] = items[i].(*model.GraphItem)
		}

		lock.Acquire()
		go func(addr string, graphItems []*model.GraphItem, count int) {
			var (
				resp   *model.RpcResponse
				err    error
				sendOk bool
			)

			defer lock.Release()

			resp = &model.RpcResponse{}
			sendOk = false

			for i := 0; i < 3; i++ { //最多重试3次
				err = GraphConnPools.Call(addr, "Graph.Send", graphItems, resp)
				if err == nil {
					sendOk = true
					break
				}
				time.Sleep(time.Millisecond * 10)
			}

			// statistics
			if !sendOk {
				dlog.Infof("send to graph %s:%s fail: %v", node, addr, err)
				proc.SendToGraphFailCnt.IncrBy(int64(count))
			} else {
				proc.SendToGraphCnt.IncrBy(int64(count))
			}
		}(addr, graphItems, count)
	}
}

// Ts db定时任务, 将数据通过api发送到ts db
func forward2TsDBTask(concurrent int) {
	var (
		batch int
		retry int
		lock  *utils.Semaphore
	)

	batch = g.Conf().TsDB.Batch // 一次发送,最多batch条数据
	retry = g.Conf().TsDB.Retry
	lock = utils.NewSemaphore(concurrent)

	for {
		var (
			items []interface{}
		)

		items = TsDBQueue.PopBackBy(batch)
		if len(items) == 0 {
			time.Sleep(DefaultSendTaskSleepInterval)
			continue
		}

		//  同步Call + 有限并发 进行发送
		lock.Acquire()
		go func(itemList []interface{}) {
			var (
				newBuffer bytes.Buffer
				count     int
				err       error
			)

			defer lock.Release()

			count = len(itemList)

			for i := 0; i < count; i++ {
				item := itemList[i].(*model.TsDBItem)
				newBuffer.WriteString(item.TsDBString())
				newBuffer.WriteString("\n")
			}

			for i := 0; i < retry; i++ {
				err = TsDBConnPools.Send(newBuffer.Bytes())
				if err == nil {
					proc.SendToTsDBCnt.IncrBy(int64(len(itemList)))
					break
				}

				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				proc.SendToTsDBFailCnt.IncrBy(int64(len(itemList)))
				dlog.Error(err)
				return
			}
		}(items)
	}
}
