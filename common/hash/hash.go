package hash

// 一致性哈希环,用于管理服务器节点.
type ConsistentHashNodeRing struct {
	ring *Consistent
}

func NewConsistentHashNodesRing(numberOfReplicas int32, nodes []string) *ConsistentHashNodeRing {
	ret := &ConsistentHashNodeRing{ring: NewConsistent()}
	ret.SetNumberOfReplicas(numberOfReplicas)
	ret.SetNodes(nodes)
	return ret
}

// 根据pk,获取node节点. hash(pk) -> node
func (u *ConsistentHashNodeRing) GetNode(pk string) (string, error) {
	return u.ring.Get(pk)
}

func (u *ConsistentHashNodeRing) SetNodes(nodes []string) {
	for _, node := range nodes {
		u.ring.Add(node)
	}
}

func (u *ConsistentHashNodeRing) SetNumberOfReplicas(num int32) {
	u.ring.NumberOfReplicas = int(num)
}
