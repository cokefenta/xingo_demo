package core

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
	"xingo_demo/pb"
	"sync"
)

type WorldMgr struct {
	PlayerNumGen int32
	Players      map[string]*Player
	sync.RWMutex
}

var WorldMgrObj *WorldMgr

func init() {
	WorldMgrObj = &WorldMgr{
		PlayerNumGen:    0,
		Players:         make(map[string]*Player),
	}
}

func (this *WorldMgr)AddPlayer(fconn iface.Iconnection) (error) {
	this.Lock()
	defer this.Unlock()
	this.PlayerNumGen += 1
	p := &NewPlayer(fconn, this.PlayerNumGen)
	this.Players[p.Pid] = p
	//出现在出生点
	this.Move(p)
	return nil
}

func (this *WorldMgr)Move(p *Player){
	data := &pb.BroadCast{
		Pid : p.Pid,
		Tp: 2,
		P: &pb.Position{
			X: p.X,
			Y: p.Y,
		},
	}
	this.Broadcast(200, data)
}

func (this *WorldMgr) GetPlayer(pid int32)(*Player, error){
	p, ok := this.Players[pid]
	if ok{
		return p, nil
	}else{
		return nil, errors.new("no player in the world!!!")
	}
}

func (this *WorldMgr) Broadcast(msgId uint32, data proto.Message) {
	for _, p := range this.Players {
		p.SendMsg(msgId, data)
	}
}