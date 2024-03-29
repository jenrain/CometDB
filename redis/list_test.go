package redis

import (
	"CometDB"
	"CometDB/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRedisObject_LPop(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-lpop")
	opts.DirPath = dir
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	res, err := rds.LPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), res)
	res, err = rds.LPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), res)
	res, err = rds.LPush(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(3), res)

	val, err := rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
}

func TestRedisObject_RPop(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-rpop")
	opts.DirPath = dir
	t.Log(dir)
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	res, err := rds.RPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), res)
	res, err = rds.RPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), res)
	res, err = rds.RPush(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(3), res)

	val, err := rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
}

func TestRedisObject_LLen(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-llen")
	opts.DirPath = dir
	t.Log(dir)
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	res, err := rds.RPush([]byte("foo1"), []byte("bar1"))
	assert.Nil(t, err)
	t.Log(res)
	res, err = rds.RPush([]byte("foo1"), []byte("bar2"))
	assert.Nil(t, err)
	t.Log(res)
	res, err = rds.RPush([]byte("foo1"), []byte("bar3"))
	assert.Nil(t, err)
	t.Log(res)

	size, err := rds.LLen([]byte("foo1"))
	assert.Nil(t, err)
	t.Log(size)
}
