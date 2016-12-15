# xingo_dem
//消息对应关系如下
/*
msgId            client                 server               描述
1                  -                    SyncPid              同步玩家本次登录的ID(用来标识玩家)
2                  Talk                   -                  世界聊天
3                  Position                   -              移动
200                -                    BroadCast            广播消息(Tp 1 世界聊天 2 移动)
201                -                    -                    广播消息 掉线
*/

sudo protoc3 -I=/home/huangxin/workspace/go_fighting/src/xingo_demo/pb --go_out=/home/huangxin/workspace/go_fighting/src/xingo_demo/pb /home/huangxin/workspace/go_fighting/src/xingo_demo/pb/*.proto
