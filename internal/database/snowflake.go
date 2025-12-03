package database

import (
	"hash/fnv"
	"net"
	"os"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
	snowflakeOnce sync.Once
	snowflakeErr  error
)

// getNodeID 自动生成节点 ID (基于 MAC 或 hostname)
// 返回值范围: 0-1023
func getNodeID() int64 {
	// 尝试从 MAC 地址生成
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				h := fnv.New64a()
				_, _ = h.Write(iface.HardwareAddr)
				return int64(h.Sum64() % 1024)
			}
		}
	}
	// 回退到 hostname
	hostname, _ := os.Hostname()
	h := fnv.New64a()
	_, _ = h.Write([]byte(hostname))
	return int64(h.Sum64() % 1024)
}

// InitSnowflake 初始化雪花算法节点
// 节点 ID 基于机器 MAC 地址或 hostname 自动生成
func InitSnowflake() error {
	snowflakeOnce.Do(func() {
		nodeID := getNodeID()
		snowflakeNode, snowflakeErr = snowflake.NewNode(nodeID)
	})
	return snowflakeErr
}

// GenerateID 生成雪花 ID
// 如果未初始化，会自动初始化
func GenerateID() int64 {
	if snowflakeNode == nil {
		_ = InitSnowflake()
	}
	return snowflakeNode.Generate().Int64()
}

// GetNodeID 获取当前节点 ID (用于调试)
func GetNodeID() int64 {
	return getNodeID()
}
