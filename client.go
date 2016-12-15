package main

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"encoding/binary"
	"bytes"
	"os"
	"os/signal"
	"xingo_demo/pb"
	"math/rand"
	"time"
)

type PkgData struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

type TcpClient struct{
	conn *net.TCPConn
	addr *net.TCPAddr
	X float32
	Y float32
	Pid int32
}

func NewTcpClient(ip string, port int) *TcpClient{
	addr := &net.TCPAddr{
		IP: net.ParseIP(ip),
		Port: port,
		Zone: "",
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err == nil{
		client := &TcpClient{
		conn: conn,
		addr: addr,
		}
		client.ConnectionMade()
		return client
	}else{
		panic(err)
	}

}

func (this *TcpClient)ConnectionMade(){
	fmt.Println("链接建立")
}

func (this *TcpClient)ConnectionLost(){
	fmt.Println("链接断开")
}

func (this *TcpClient) Unpack(headdata []byte) (head *PkgData, err error) {
	headbuf := bytes.NewReader(headdata)

	head = &PkgData{}

	// 读取Len
	if err = binary.Read(headbuf, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	// 读取MsgId
	if err = binary.Read(headbuf, binary.LittleEndian, &head.MsgId); err != nil {
		return nil, err
	}

	// 封包太大
	//if head.Len > MaxPacketSize {
	//	return nil, packageTooBig
	//}

	return head, nil
}

func (this *TcpClient) Pack(msgId uint32, data proto.Message) (out []byte, err error) {
	outbuff := bytes.NewBuffer([]byte{})
	// 进行编码
	dataBytes := []byte{}
	if data != nil {
		dataBytes, err = proto.Marshal(data)
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
	}
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return
	}
	// 写MsgId
	if err = binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
		return
	}

	//all pkg data
	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return
	}

	out = outbuff.Bytes()
	return

}

func (this *TcpClient)DoMsg(pdata *PkgData){
	//处理消息
	fmt.Println(fmt.Sprintf("msg id :%d, data len: %d", pdata.MsgId, pdata.Len))
	if pdata.MsgId == 1{
		syncpid := &pb.SyncPid{}
		proto.Unmarshal(pdata.Data, syncpid)
		this.Pid = syncpid.Pid
	}else if pdata.MsgId == 200{
		bdata := &pb.BroadCast{}
		proto.Unmarshal(pdata.Data, bdata)
		if bdata.Tp == 2{
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
		}else{
			fmt.Println(fmt.Sprintf("世界聊天,%i: %s", bdata.Pid, bdata.GetContent()))
		}
		//聊天或者移动
		time.Sleep(3*time.Second)
		tp := rand.Intn(2)
		if tp == 0{
			//聊天
			msg := &pb.Talk{
				Content: "你猜猜我是谁？",
			}
			this.Send(2, msg)
		}else{
			//移动
			msg := &pb.Position{
				X: this.X + 1,
				Y : this.Y + 1,
			}
			this.Send(3, msg)
		}
	}
}

func (this *TcpClient)Send(msgID uint32, data proto.Message){
	dd, err := this.Pack(msgID, data)
	if err == nil{
		this.conn.Write(dd)
	}else{
		fmt.Println(err)
	}

}

func (this *TcpClient)Start(){
	go func() {
		for {
		//read per head data
		headdata := make([]byte, 8)

		if _, err := io.ReadFull(this.conn, headdata); err != nil {
			fmt.Println(err)
			this.ConnectionLost()
			return
		}
		pkgHead, err := this.Unpack(headdata)
		if err != nil {
			this.ConnectionLost()
			return
		}
		//data
		if pkgHead.Len > 0 {
			pkgHead.Data = make([]byte, pkgHead.Len)
			if _, err := io.ReadFull(this.conn, pkgHead.Data); err != nil {
				this.ConnectionLost()
				return
			}
		}
		this.DoMsg(pkgHead)
	}
	}()
}

func main() {
	for i := 0; i< 5; i ++{
		client := NewTcpClient("0.0.0.0", 8909)
		client.Start()
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
}