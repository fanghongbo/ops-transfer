package sender

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/proc"
	"github.com/fanghongbo/ops-transfer/utils"
	"time"
)

const (
	DefaultProcCronPeriod = time.Duration(5) * time.Second    //ProcCron的周期,默认1s
	DefaultLogCronPeriod  = time.Duration(3600) * time.Second //LogCron的周期,默认300s
)

func startSenderCron() {
	go startProcCron()
	go startLogCron()
}

func startProcCron() {
	for {
		time.Sleep(DefaultProcCronPeriod)
		refreshSendingCacheSize()
	}
}

func startLogCron() {
	for {
		time.Sleep(DefaultLogCronPeriod)
		logConnPoolsProc()
	}
}

func refreshSendingCacheSize() {
	var cfg *g.GlobalConfig

	proc.JudgeQueuesCnt.SetCnt(calcSendCacheSize(JudgeQueues))
	proc.GraphQueuesCnt.SetCnt(calcSendCacheSize(GraphQueues))

	cfg = g.Conf()

	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		proc.TsDBQueuesCnt.SetCnt(int64(TsDBQueue.Len()))
	}
}

func calcSendCacheSize(mapList map[string]*utils.ListLimited) int64 {
	var cnt int64 = 0

	for _, list := range mapList {
		if list != nil {
			cnt += int64(list.Len())
		}
	}
	return cnt
}

func logConnPoolsProc() {
	dlog.Infof("judge connPools proc: %v", JudgeConnPools.Proc())
	dlog.Infof("graph connPools proc: %v", GraphConnPools.Proc())
}
