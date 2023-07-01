package main

import (
	"CometDB"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func handlePut(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		log.Printf("[handlerPut] fail to get body data||reason:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to get body data",
			"code":    http.StatusBadRequest,
		})
		return
	}
	var body map[string]string
	err = json.Unmarshal(data, &body)
	if err != nil {
		log.Printf("[handlerPut] fail to parse body data||reason:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail to parse body data",
			"code":    http.StatusBadRequest,
		})
		return
	}
	for k, v := range body {
		if err := db.Put([]byte(k), []byte(v)); err != nil {
			log.Printf("[handlerPut] fail to put data||key=%v||value=%v||reason:%v\n", k, v, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("fail to put data||key=%v||value=%v", k, v),
				"code":    http.StatusInternalServerError,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "nil",
		"code":    http.StatusOK,
	})
}

func handleGet(c *gin.Context) {
	key := c.Query("key")
	v, err := db.Get([]byte(key))
	if err != nil && err != CometDB.ErrKeyNotFound {
		log.Printf("[handleGet] fail to get value||key=%v||reason:%v\n", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fail to get data",
			"code":    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"value": string(v),
	})
}

func handleDelete(c *gin.Context) {
	key := c.Query("key")
	err := db.Delete([]byte(key))
	if err != nil {
		log.Printf("[handleGet] fail to delete data||key=%v||reason:%v\n", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("fail to delete data||key=%v", key),
			"code":    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "nil",
	})
}

func handleListKeys(c *gin.Context) {
	keys := db.ListKeys()
	var res []string
	for _, key := range keys {
		fmt.Println("key:", string(key))
		res = append(res, string(key))
	}
	c.JSON(http.StatusOK, res)
}

func handleStat(c *gin.Context) {
	stat := db.Stat()
	c.JSON(http.StatusOK, stat)
}
