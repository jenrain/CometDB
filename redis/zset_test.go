package redis

import (
	"CometDB"
	"CometDB/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRedisObject_ZScore(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-zset")
	opts.DirPath = dir
	rds, err := NewRedisObject(opts)
	assert.Nil(t, err)

	ok, err := rds.ZAdd(utils.GetTestKey(1), 113, []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.ZAdd(utils.GetTestKey(1), 333, []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.ZAdd(utils.GetTestKey(1), 98, []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	score, err := rds.ZScore(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, float64(333), score)
	score, err = rds.ZScore(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.Equal(t, float64(98), score)
}
