package apis

import (
	"MMO-HappyRunning/core"
	"MMO-HappyRunning/pb"
	"fmt"
	"github.com/vastea/myzinx/ziface"
	"github.com/vastea/myzinx/znet"
	"google.golang.org/protobuf/proto"
)

type WorldChat struct {
	znet.BaseRouter
}

func (wc *WorldChat) Handle(request ziface.IRequest) {
	// 1-解析客户端传递进来的proto协议
	protoMsg := &pb.Talk{}
	dataBytes := request.GetData()
	err := proto.Unmarshal(dataBytes, protoMsg)
	if err != nil {
		fmt.Println("[ERROR] Unmarshal error:", err)
		return
	}
	// 2-当前的聊天数据是属于哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")
	// 3-根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 4-将这个消息广播给其他全部在线的玩家
	player.Talk(protoMsg.Content)
}
