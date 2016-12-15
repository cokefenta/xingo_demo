package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
	"xingo_demo/pb"
	"math/rand"
)

type Player struct {
	Fconn iface.Iconnection
	Pid   int32
	X     float32
	Y     float32
}

func NewPlayer(fconn iface.Iconnection, pid int32) *Player {
	p := &Player{
		Fconn: fconn,
		Pid:   pid,
		X:     float32(rand.Intn(161) + 160),
		Y:     float32(rand.Intn(217) + 134),
	}
	return p
}

func (this *Player) UpdatePos(x float32, y float32) {
	this.X = x
	this.Y = y
	WorldMgrObj.Move(this)
}

func (this *Player) Talk(content string){
	data := &pb.BroadCast{
		Pid : this.Pid,
		Tp: 1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	WorldMgrObj.Broadcast(200, data)
}

func (this *Player) LostConnection(){
	msg := &pb.SyncPid{
		Pid: this.Pid,
	}
	WorldMgrObj.Broadcast(201, msg)
}

func (this *Player) SendMsg(msgId uint32, data proto.Message) {
	if this.Fconn != nil {
		this.Fconn.Send(msgId, data)
	}
}
