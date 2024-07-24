// Package core 核心算法包
package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子类型
type Grid struct {
	// 格子ID
	GID int
	// 格子的左边边界坐标
	MinX int
	// 格子的右边边界坐标
	MaxX int
	// 格子的上边边界坐标
	MinY int
	// 格子的下边边界坐标
	MaxY int
	// 当前格子内玩家或物体的ID集合
	PlayerIDs map[int]bool
	// 保护当前集合的锁
	PIDLock sync.RWMutex
}

// NewGrid 初始化当前格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		PlayerIDs: make(map[int]bool),
	}
}

// AddPlayer 给格子内添加一个玩家
func (g *Grid) AddPlayer(PlayerID int) {
	g.PIDLock.Lock()
	defer g.PIDLock.Unlock()

	g.PlayerIDs[PlayerID] = true
}

// RemovePlayer 从格子中删除一个玩家
func (g *Grid) RemovePlayer(PlayerID int) {
	g.PIDLock.Lock()
	defer g.PIDLock.Unlock()

	delete(g.PlayerIDs, PlayerID)
}

// GetPlayerIDs 得到当前格子中所有的玩家/物品
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.PIDLock.RLock()
	defer g.PIDLock.RUnlock()

	for playerID, _ := range g.PlayerIDs {
		playerIDs = append(playerIDs, playerID)
	}
	return
}

// String 重写String()方法，使用fmt.Println打印Grid时，底层会默认调用String()方法，打印出格子的基本信息(调试用)
func (g *Grid) String() string {
	return fmt.Sprintf("==========[Grid--%d]==========\n"+
		"minX: %d, maxX: %d, minY: %d, maxY: %d\n"+
		"playerIDs: %v\n"+
		"==========[Grid--%d]==========\n",
		g.GID, g.MinX, g.MaxY, g.MinY, g.MaxY, g.PlayerIDs, g.GID)
}
