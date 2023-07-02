package main

import (
	"CometDB"
	"CometDB/redis"
	"github.com/tidwall/redcon"
	"log"
	"sync"
)

const addr = "127.0.0.1:6380"

type CometDBServer struct {
	dbs    *redis.RedisObject
	server *redcon.Server
	mu     sync.RWMutex
}

func main() {
	// 打开 Redis 数据结构服务
	redisObject, err := redis.NewRedisObject(CometDB.DefaultOptions)
	if err != nil {
		panic(err)
	}

	// 初始化服务端
	cometDBServer := &CometDBServer{
		dbs: &redis.RedisObject{},
	}
	cometDBServer.dbs = redisObject
	cometDBServer.server = redcon.NewServer(addr, execClientCommand, cometDBServer.accept, cometDBServer.close)
	cometDBServer.listen()
}

func (svr *CometDBServer) listen() {
	log.Println("cometDB server is running, ready to accept connections.")
	_ = svr.server.ListenAndServe()
}

func (svr *CometDBServer) accept(conn redcon.Conn) bool {
	cli := new(CometDBClient)
	svr.mu.Lock()
	defer svr.mu.Unlock()

	cli.server = svr
	cli.db = svr.dbs
	conn.SetContext(cli)
	return true
}

func (svr *CometDBServer) close(conn redcon.Conn, err error) {
	_ = svr.dbs.Close()
	_ = svr.server.Close()
}
