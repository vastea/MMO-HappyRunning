package main

import (
	"MMO-HappyRunning/apis"
	"MMO-HappyRunning/core"
	"fmt"
	"github.com/vastea/myzinx/ziface"
	"github.com/vastea/myzinx/znet"
)

func main() {
	// 创建myzinx server
	server := znet.NewServer()
	// 链接创建和销毁的hook函数
	server.SetOnConnectionStart(playerOnline)
	// 注册路由
	server.AddRouter(2, &apis.WorldChat{})
	// 启动服务
	server.Serve()
}

// OnConnectionAdd 当前客户端建立连接之后的HOOK函数
func playerOnline(connection ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(connection)
	// 给客户端发送MsgID为1的消息，同步当前player的ID给客户端
	player.SyncPID()
	// 给客户端发送MsgID为200的消息，同步当前player的初始位置给客户端
	player.BroadCastStartPosition()
	// 将当前新上线的玩家添加到world中
	core.WorldMgrObj.AddPlayer(player)
	// 将pid和connection绑定
	player.Connection.SetProperty("pid", player.Pid)

	fmt.Println("[PLAYER-ONLINE] Player has online, the Pid is:", player.Pid)
}
