package sender

import (
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/utils"
)

const (
	DefaultSendQueueMaxSize = 102400 //10.24w
)

// 发送缓存队列
// node -> queue_of_data
var (
	TsDBQueue   *utils.ListLimited
	JudgeQueues = make(map[string]*utils.ListLimited)
	GraphQueues = make(map[string]*utils.ListLimited)
)

func initSendQueues() {
	var cfg *g.GlobalConfig

	cfg = g.Conf()

	for node := range cfg.Judge.Cluster {
		Q := utils.NewSafeListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	for node, item := range g.GraphCluster {
		for _, addr := range item.Addrs {
			Q := utils.NewSafeListLimited(DefaultSendQueueMaxSize)
			GraphQueues[node+addr] = Q
		}
	}

	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		TsDBQueue = utils.NewSafeListLimited(DefaultSendQueueMaxSize)
	}
}
