package api

import (
	"xingo_demo/pb"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/fnet"
	"github.com/viphxin/xingo/logger"
	_ "time"
	"fmt"
)

type TestRouter struct {
}

/*
ping test
*/
func (this *TestRouter) Api_0(request *fnet.PkgAll) {
	logger.Debug("call Api_0")
	// request.Fconn.SendBuff(0, nil)
	request.Fconn.Send(0, nil)
}

/*
世界聊天
 */
func (this *TestRouter) Api_2(request *fnet.PkgAll) {

}

/*
移动
 */
func (this *TestRouter) Api_3(request *fnet.PkgAll) {
	msg := &pb.Position{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user move: x: %s y: %s", userId, tocken))
		if userId != "" {
			request.Fconn.SetProperty("uid", userId)
			resp := &pb.CommonResponse{
				State: 1,
			}
			request.Fconn.SendBuff(1, resp)
		} else {
			logger.Error("no userid found")
			request.Fconn.LostConnection()
		}
	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}