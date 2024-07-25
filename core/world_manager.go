package core

import "sync"

const (
	AOI_MIN_X int = 85
	AOI_MAX_X int = 410
	AOI_CNT_X int = 10
	AOI_MIN_Y int = 75
	AOI_MAX_Y int = 400
	AOI_CNT_Y int = 20
)

type WorldManager struct {
	// AOIManager 当前世界地图AOI的管理模块
	aoiMgr *AOIManager
	// 当前全部在线的Player集合
	players map[int32]*Player
	// 保护Player集合的锁
	pLock sync.RWMutex
}

var WorldMgrObj *WorldManager

// AddPlayer 添加一个在线玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.players[player.Pid] = player
	wm.pLock.Unlock()
	// 将player添加到AOIManager中
	wm.aoiMgr.AddPID2GIDByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 删除一个在线玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	player := wm.players[pid]
	wm.pLock.Lock()
	delete(wm.players, pid)
	wm.pLock.Unlock()
	wm.aoiMgr.RemovePIDFromGIDByPos(int(player.Pid), player.X, player.Z)
}

// GetPlayerByPid 通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	player := wm.players[pid]
	wm.pLock.RUnlock()
	return player
}

// GetAllPlayers 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)
	for _, player := range wm.players {
		players = append(players, player)
	}
	return players
}

// GetAoiManager 获取AOI地图管理模块
func (wm *WorldManager) GetAoiManager() *AOIManager {
	return wm.aoiMgr
}

func init() {
	WorldMgrObj = &WorldManager{
		// 创建世界AOI地图规划
		aoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
		// 初始化players集合
		players: make(map[int32]*Player),
	}
}
