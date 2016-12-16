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
	Players      map[int32]*Player
	sync.RWMutex
}

var WorldMgrObj *WorldMgr

func init() {
	WorldMgrObj = &WorldMgr{
		PlayerNumGen:    0,
		Players:         make(map[int32]*Player),
	}
}

func (this *WorldMgr)AddPlayer(fconn iface.Iconnection) (*Player, error) {
	this.Lock()
	defer this.Unlock()
	this.PlayerNumGen += 1
	p := NewPlayer(fconn, this.PlayerNumGen)
	this.Players[p.Pid] = p
	//同步Pid
	msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, msg)
	//出现在出生点
	this.Move(p, -1)
	return p, nil
}

func (this *WorldMgr)RemovePlayer(pid int32){
	this.Lock()
	defer this.Unlock()
	delete(this.Players, pid)
}

func (this *WorldMgr)Move(p *Player, action int32){
	var data *pb.BroadCast
	if action == -1{
		//出生
		data = &pb.BroadCast{
			Pid : p.Pid,
			Tp: 2,
			Data: &pb.BroadCast_P{
				P: &pb.Position{
				X: p.X,
				Y: p.Y,
				},
			},
		}
	}else{
		//不广播坐标, 广播动作数据
		data = &pb.BroadCast{
			Pid : p.Pid,
			Tp: 3,
			Data: &pb.BroadCast_ActionData{
				ActionData: action,
			},
		}
	}
	this.Broadcast(200, data)
}

func (this *WorldMgr) GetPlayer(pid int32)(*Player, error){
	p, ok := this.Players[pid]
	if ok{
		return p, nil
	}else{
		return nil, errors.New("no player in the world!!!")
	}
}

func (this *WorldMgr) Broadcast(msgId uint32, data proto.Message) {
	for _, p := range this.Players {
		p.SendMsg(msgId, data)
	}
}