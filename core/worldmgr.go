package core

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
	"xingo_demo/pb"
	"github.com/viphxin/xingo/logger"
)

type WorldMgr struct {
	PlayerNumGen int32
	Players      map[int32]*Player
	AoiObj1       *AOIMgr//地图1
}

var WorldMgrObj *WorldMgr

func init() {
	WorldMgrObj = &WorldMgr{
		PlayerNumGen:    0,
		Players:         make(map[int32]*Player),
		AoiObj1:          NewAOIMgr(85, 410, 75, 400, 10, 20),
	}
}

func (this *WorldMgr)AddPlayer(fconn iface.Iconnection) (*Player, error) {
	this.PlayerNumGen += 1
	p := NewPlayer(fconn, this.PlayerNumGen)
	this.Players[p.Pid] = p
	//同步Pid
	msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, msg)
	//加到aoi
	this.AoiObj1.Add2AOI(p)
	//周围的人
	p.SyncSurrouding()
	return p, nil
}

func (this *WorldMgr)RemovePlayer(pid int32){
	//从aoi移除
	this.AoiObj1.LeaveAOI(this.Players[pid])
	delete(this.Players, pid)
}

func (this *WorldMgr)Move(p *Player){
	var data *pb.BroadCast
	data = &pb.BroadCast{
		Pid : p.Pid,
		Tp: 4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
			},
		},
	}
	/*aoi*/
	pids, err := this.AoiObj1.GetSurroundingPids(p)
	if err == nil{
		for _, pid := range pids{
			player, err1 := this.GetPlayer(pid)
			if err1 == nil{
				player.SendMsg(200, data)
			}
		}
	}
}

func (this *WorldMgr)SendMsgByPid(pid int32, msgId uint32, data proto.Message){
	p, err := this.GetPlayer(pid)
	if err == nil{
		p.SendMsg(msgId, data)
	}
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

func (this *WorldMgr) BroadcastBuff(msgId uint32, data proto.Message) {
	for _, p := range this.Players {
		p.SendBuffMsg(msgId, data)
	}
}

func (this *WorldMgr) AOIBroadcast(p *Player, msgId uint32, data proto.Message) {
	/*aoi*/
	pids, err := WorldMgrObj.AoiObj1.GetSurroundingPids(p)
	if err == nil{
		for _, pid := range pids{
			player, err1 := WorldMgrObj.GetPlayer(pid)
			if err1 == nil {
				player.SendMsg(msgId, data)
			}
		}
	}else{
		logger.Error(err)
	}
}