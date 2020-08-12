package sender

import (
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/pool"
	"github.com/fanghongbo/ops-transfer/utils"
)

// 连接池
// node_address -> connection_pool
var (
	JudgeConnPools *pool.RpcConnPools
	TsDBConnPools  *pool.TsDBConnPools
	GraphConnPools *pool.RpcConnPools
)

func initConnPools() {
	var (
		cfg            *g.GlobalConfig
		judgeInstances *utils.StringSet
		graphInstances *utils.SafeSet
	)

	cfg = g.Conf()

	// judge
	judgeInstances = utils.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}

	JudgeConnPools = pool.CreateRpcConnPools(
		cfg.Judge.MaxConn,
		cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout,
		cfg.Judge.CallTimeout,
		judgeInstances.ToSlice(),
	)

	// graph
	graphInstances = utils.NewSafeSet()
	for _, item := range g.GraphCluster {
		for _, addr := range item.Addrs {
			graphInstances.Add(addr)
		}
	}

	GraphConnPools = pool.CreateRpcConnPools(
		cfg.Graph.MaxConn,
		cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout,
		cfg.Graph.CallTimeout,
		graphInstances.ToSlice(),
	)

	// ts db
	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		TsDBConnPools = pool.NewTsDBConnPool(
			cfg.TsDB.Addr,
			cfg.TsDB.MaxConn,
			cfg.TsDB.MaxIdle,
			cfg.TsDB.ConnTimeout,
			cfg.TsDB.CallTimeout,
		)
	}
}

func DestroyConnPools() {
	var cfg *g.GlobalConfig

	cfg = g.Conf()

	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()

	if cfg.TsDB != nil && cfg.TsDB.Enabled {
		TsDBConnPools.Destroy()
	}
}
