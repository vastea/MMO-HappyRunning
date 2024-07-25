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
