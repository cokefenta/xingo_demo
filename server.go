package main

import (
	"github.com/viphxin/xingo/fserver"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"xingo_demo/api"
	"xingo_demo/core"

	_ "net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
	_ "time"
)

func DoConnectionMade(fconn iface.Iconnection) {
	logger.Debug("111111111111111111111111")
	p, _ := core.WorldMgrObj.AddPlayer(fconn)
	fconn.SetProperty("pid", p.Pid)
}

func DoConnectionLost(fconn iface.Iconnection) {
	logger.Debug("222222222222222222222222")
	pid, _ := fconn.GetProperty("pid")
	p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
	//移除玩家
	core.WorldMgrObj.RemovePlayer(pid.(int32))
	//消失在地图
	p.LostConnection()
}

func DoStop(){
	logger.Debug("onstop !!!!!!!!!!!!!!")
}

func main() {
	s := fserver.NewServer()

	//add api ---------------start
	TestRouterObj := &api.TestRouter{}
	s.AddRouter(TestRouterObj)
	//add api ---------------end
	//regest callback
	utils.GlobalObject.OnConnectioned = DoConnectionMade
	utils.GlobalObject.OnClosed = DoConnectionLost
	utils.GlobalObject.OnServerStop = DoStop

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6061", nil))
	// 	// for {
	// 	// 	time.Sleep(time.Second * 10)
	// 	// 	fm, err := os.OpenFile("./memory.log", os.O_RDWR|os.O_CREATE, 0644)
	// 	// 	if err != nil {
	// 	// 		fmt.Println(err)
	// 	// 	}
	// 	// 	pprof.WriteHeapProfile(fm)
	// 	// 	fm.Close()
	// 	// }
	// }()
	s.Serve()
}
