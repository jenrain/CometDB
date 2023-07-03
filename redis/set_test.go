package redis

import (
	"CometDB"
	"CometDB/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRedisObject_SIsMember(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-sismember")
	opts.DirPath = dir
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	ok, err := rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SIsMember(utils.GetTestKey(2), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-not-exist"))
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRedisObject_SRem(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-srem")
	opts.DirPath = dir
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	ok, err := rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SRem(utils.GetTestKey(2), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SRem(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.False(t, ok)
}
