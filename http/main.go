package main

import (
	"CometDB"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var db *CometDB.DB

func init() {
	// 初始化 DB 实例
	var err error
	options := CometDB.DefaultOptions
	dir, err := os.MkdirTemp("", "cometDB-http")
	log.Printf("db dir: %v\n", dir)
	options.DirPath = dir
	db, err = CometDB.Open(options)
	if err != nil {
		panic(fmt.Sprintf("failed to open db: %v", err))
	}
}

func main() {
	r := gin.Default()
	initRouter(r)
	r.Run(":8080")
}
