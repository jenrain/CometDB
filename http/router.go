package main

import "github.com/gin-gonic/gin"

func initRouter(r *gin.Engine) {
	r.PUT("/put", handlePut)
	r.GET("/get", handleGet)
	r.DELETE("/del", handleDelete)
	r.GET("/keys", handleListKeys)
	r.GET("/stat", handleStat)
}
