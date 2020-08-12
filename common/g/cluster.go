package g

import "strings"

var (
	GraphCluster map[string]*ClusterNode
)

type ClusterNode struct {
	Addrs []string `json:"addrs"`
}

func NewClusterNode(addrs []string) *ClusterNode {
	return &ClusterNode{addrs}
}

// map["node"]="host1,host2" --> map["node"]=["host1", "host2"]
func formatClusterItems(cluster map[string]string) map[string]*ClusterNode {
	var data map[string]*ClusterNode

	data = make(map[string]*ClusterNode)
	for node, clusterStr := range cluster {
		var (
			items []string
			temp  []string
		)

		items = strings.Split(clusterStr, ",")
		temp = make([]string, 0)
		for _, item := range items {
			temp = append(temp, strings.TrimSpace(item))
		}
		data[node] = NewClusterNode(temp)
	}

	return data
}

// 初始化集群配置
func InitClusterNode() {
	GraphCluster = formatClusterItems(config.Graph.Cluster)
}
