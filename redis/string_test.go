package redis

import (
	"CometDB"
	"CometDB/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRedisObject_Get(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-get")
	opts.DirPath = dir
	r, err := NewRedisObject(opts)
	assert.Nil(t, err)

	err = r.Set([]byte("foo1"), 0, []byte("bar1"))
	assert.Nil(t, err)
	err = r.Set([]byte("foo2"), time.Second*5, []byte("bar2"))
	assert.Nil(t, err)

	val1, err := r.Get([]byte("foo1"))
	assert.Nil(t, err)
	assert.NotNil(t, val1)
	t.Log(string(val1))

	val2, err := r.Get([]byte("foo2"))
	assert.Nil(t, err)
	assert.NotNil(t, val2)
	t.Log(string(val2))

	val3, err := r.Get([]byte("foo3"))
	assert.Nil(t, err)
	assert.NotNil(t, val3)
	t.Log(string(val3))
}

func TestRedisObject_Del_Type(t *testing.T) {
	opts := CometDB.DefaultOptions
	dir, _ := os.MkdirTemp("", "CometDB-redis-del-type")
	opts.DirPath = dir
	r, err := NewRedisObject(opts)
	assert.Nil(t, err)

	// del
	err = r.Del(utils.GetTestKey(11))
	assert.Nil(t, err)

	err = r.Set(utils.GetTestKey(1), 0, utils.RandomValue(100))
	assert.Nil(t, err)

	// type
	typ, err := r.Type(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.Equal(t, String, typ)

	err = r.Del(utils.GetTestKey(1))
	assert.Nil(t, err)

	_, err = r.Get(utils.GetTestKey(1))
	assert.Equal(t, CometDB.ErrKeyNotFound, err)
}
