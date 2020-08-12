package sender

import (
	"github.com/fanghongbo/ops-transfer/common/g"
	"github.com/fanghongbo/ops-transfer/common/hash"
	"sort"
)

// 服务节点的一致性哈希环
// pk -> node
var (
	JudgeNodeRing *hash.ConsistentHashNodeRing
	GraphNodeRing *hash.ConsistentHashNodeRing
)

func initNodeRings() {
	var cfg *g.GlobalConfig

	cfg = g.Conf()

	JudgeNodeRing = hash.NewConsistentHashNodesRing(int32(cfg.Judge.Replicas), getClusterNodeName(cfg.Judge.Cluster))
	GraphNodeRing = hash.NewConsistentHashNodesRing(int32(cfg.Graph.Replicas), getClusterNodeName(cfg.Graph.Cluster))
}

func getClusterNodeName(m map[string]string) []string {
	var (
		keys sort.StringSlice
		i    int
	)

	keys = make(sort.StringSlice, len(m))
	i = 0
	for key := range m {
		keys[i] = key
		i++
	}

	keys.Sort()
	return []string(keys)
}
