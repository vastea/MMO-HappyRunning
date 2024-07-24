// Package core 核心算法包
package core

import "fmt"

// AOIManager AOI区域管理模块
type AOIManager struct {
	MinX  int           //区域的左边界坐标
	MaxX  int           //区域的右边界坐标
	CntX  int           //X轴方向格子的数量
	MinY  int           //区域的上边界坐标
	MaxY  int           //区域的下边界坐标
	CntY  int           //Y轴方向格子的数量
	grids map[int]*Grid //当前区域中有哪些格子map -- key是格子的id，value是格子的对象
}

// NewAOIManager 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY int) *AOIManager {
	aoim := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntX:  cntX,
		MinY:  minY,
		MaxY:  maxY,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域的所有格子进行编号和初始化
	for y := 0; y < cntX; y++ {
		for x := 0; x < cntY; x++ {
			// 根据x、y轴上的编号(非坐标，比如x轴上有n个格子，指的是某一格子位于x轴上的第几个)，计算格子ID
			gid := y*cntX + x
			// 初始化格子
			gridWidth := aoim.gridWidth()
			gridHeight := aoim.gridHeight()
			aoim.grids[gid] = NewGrid(gid,
				aoim.MinX+gridWidth*x, aoim.MinX+gridWidth*(x+1),
				aoim.MinY+gridHeight*y, aoim.MinY+gridHeight*(y+1))
		}
	}

	return aoim
}

// GetSurroundGridsByGid 根据GID得到周边九宫格的grid集合
func (aoim *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := aoim.grids[gID]; !ok {
		return
	}
	// 将当前gID对应的grid加入到切片中
	grids = append(grids, aoim.grids[gID])
	// 判断当前grid的左右两侧是否有grid
	// 获取当前grid在X轴上的编号
	idX := gID % aoim.CntX
	// 判断当前grid的idx编号左边是否还有格子，如果有，放入gIDs集合中
	if idX > 0 {
		grids = append(grids, aoim.grids[gID-1])
	}
	// 判断当前grid的idx编号右边是否还有格子，如果有，放入gIDs集合中
	if idX < aoim.CntX-1 {
		grids = append(grids, aoim.grids[gID+1])
	}
	// 将X轴当前的格子都取出，进行遍历，判断每个格子的上下侧是否还有格子
	gIDs := make([]int, 0, len(grids))
	for _, grid := range grids {
		gIDs = append(gIDs, grid.GID)
	}
	for _, gID := range gIDs {
		// 获取当前grid在Y轴上的编号
		idY := gID / aoim.CntY
		// gid上边是否还有格子
		if idY > 0 {
			grids = append(grids, aoim.grids[gID-aoim.CntX])
		}
		if idY < aoim.CntY-1 {
			grids = append(grids, aoim.grids[gID+aoim.CntX])
		}
	}
	return
}

// GetPIDsByPos 根据X、Y坐标得到周边九宫格格子内的PlayerID集合
func (aoim *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家所在格子的GID
	gID := aoim.getGIDByPos(x, y)
	// 通过GID得到周边九宫格的信息
	grids := aoim.GetSurroundGridsByGid(gID)
	// 将九宫格里全部的playerID添加到playerIDs中
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	fmt.Printf("[DEBUG] Position x: %f, y: %f surround playerIDs are:\n"+
		"%v\n", x, y, playerIDs)
	return
}

// AddPID2GID 添加一个playerID到一个格子中
func (aoim *AOIManager) AddPID2GID(pID, gID int) {
	aoim.grids[gID].AddPlayer(pID)
}

// RemovePIDFromGID 移除一个格子中的playerID
func (aoim *AOIManager) RemovePIDFromGID(pID, gID int) {
	aoim.grids[gID].RemovePlayer(pID)
}

// AddPID2GIDByPos 通过坐标将player添加到一个格子中
func (aoim *AOIManager) AddPID2GIDByPos(pID int, x, y float32) {
	gID := aoim.getGIDByPos(x, y)
	aoim.AddPID2GID(pID, gID)
}

// RemovePIDFromGIDByPos 通过坐标把一个player从一个格子中删除
func (aoim *AOIManager) RemovePIDFromGIDByPos(pID int, x, y float32) {
	gID := aoim.getGIDByPos(x, y)
	aoim.RemovePIDFromGID(pID, gID)
}

// GetPIDsByGID 通过GID获取全部的playerID
func (aoim *AOIManager) GetPIDsByGID(gID int) (pIDs []int) {
	pIDs = aoim.grids[gID].GetPlayerIDs()
	return
}

// 打印格子信息
func (aoim *AOIManager) String() string {
	s := fmt.Sprintf("==========[AOIManager]==========\n"+
		"minX: %d, maxX: %d, cntX: %d,\n"+
		"minY: %d, maxY: %d, cntY: %d,\n"+
		"==========[AOIManager]==========\n",
		aoim.MinX, aoim.MaxX, aoim.CntX, aoim.MinY, aoim.MaxY, aoim.CntY)
	for _, grid := range aoim.grids {
		fmt.Println(grid)
	}
	return s
}

// 得到每个格子在x轴方向的宽度
func (aoim *AOIManager) gridWidth() int {
	return (aoim.MaxX - aoim.MinX) / aoim.CntX
}

// 得到每个格子的y轴方向的高度
func (aoim *AOIManager) gridHeight() int {
	return (aoim.MaxY - aoim.MinY) / aoim.CntY
}

// 通过横纵坐标获取当前所在格子的GID
func (aoim *AOIManager) getGIDByPos(x, y float32) int {
	idX := (int(x) - aoim.MinX) / aoim.gridWidth()
	idY := (int(y) - aoim.MinY) / aoim.gridHeight()
	return idX + idY*aoim.CntX
}
