package core

import (
	"MMO-HappyRunning/pb"
	"fmt"
	"github.com/vastea/myzinx/ziface"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

type Player struct {
	Pid        int32              // 玩家ID
	Connection ziface.IConnection // 当前玩家的连接(用于和客户端连接)
	X          float32            // 平面的X坐标
	Y          float32            // 高度
	Z          float32            // 平面的Y坐标
	V          float32            // 倾斜角度(0-360)
}

// 全局计数器
var pidGen int32 = 0
var pidGenLock sync.Mutex

func NewPlayer(connection ziface.IConnection) *Player {
	// 生成一个玩家ID
	pidGenLock.Lock()
	id := pidGen
	pidGen++
	pidGenLock.Unlock()
	// 创建一个玩家对象
	p := &Player{
		Pid:        id,
		Connection: connection,
		X:          float32(160 + rand.Intn(10)), // 随机在160坐标点，基于X轴若干偏移
		Y:          0,
		Z:          float32(140 + rand.Intn(20)),
		V:          0, // 角度为0
	}
	return p
}

// SendMessage 提供一个发送给客户端消息的方法，主要是将pb的protobuf数据序列化之后再调用zinx的SendMessage方法
func (p *Player) SendMessage(msgID uint32, data proto.Message) {
	// 将protoMessage结构体序列化，转换成二进制
	message, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("[ERROR] Proto Marshal error:", err)
		return
	}
	// 将二进制数据，通过myzinx框架的sendMsg将数据发送给客户端
	if p.Connection == nil {
		fmt.Println("[ERROR] Connection in player is nil")
		return
	}
	if err := p.Connection.SendMessage(msgID, message); err != nil {
		fmt.Println("[ERROR] Player sendMsg by myzinx error:", err)
		return
	}
}

// SyncPID 告知客户端玩家PID，同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {
	// 组建MsgID为1的proto消息
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 将消息发送给客户端
	p.SendMessage(1, protoMsg)
}

// BroadCastStartPosition 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID为200的proto消息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 将消息发送给客户端
	p.SendMessage(200, protoMsg)
}

// Talk 广播聊天信息
func (p *Player) Talk(content string) {
	// 1-组建一个msgID为200的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	// 2-得到当前所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()
	// 3-向所有的玩家(包括自己)发送msgID为200的消息
	for _, player := range players {
		// 每个player分别给对应的客户端发送消息
		player.SendMessage(200, protoMsg)
	}
}

// SyncSurrounding 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1-获取当前玩家周围的玩家有哪些
	pids := WorldMgrObj.aoiMgr.GetPIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	// 2-将当前玩家的位置信息通过MsgID：200发送给周围玩家(让其他玩家看到自己)
	// 2.1-组建MsgID为200的消息
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2-周围全部的玩家都向自己的客户端发送200消息
	for _, player := range players {
		player.SendMessage(200, protoMsg)
	}
	// 3-将周围的全部玩家的位置信息发送给当前的玩家客户端(让自己看到其他玩家)
	// 3.1-组建MsgID为202的消息
	playerPos := make([]*pb.Player, 0, len(pids))
	for _, player := range players {
		PlayerPosProtoMsg := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.Z,
			},
		}
		playerPos = append(playerPos, PlayerPosProtoMsg)
	}
	protoMsg202 := &pb.SyncPlayers{
		Ps: playerPos,
	}
	// 3.2-将组建好的数据发送给当前玩家的客户端
	p.SendMessage(202, protoMsg202)
}
