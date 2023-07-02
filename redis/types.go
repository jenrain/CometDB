package redis

import (
	"CometDB"
	"errors"
)

var (
	ErrWrongTypeOperation = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
)

type redisObjectType = byte

const (
	String redisObjectType = iota
	Hash
	Set
	List
	ZSet
)

type RedisObject struct {
	db *CometDB.DB
}

// NewRedisObject 初始化 Redis 数据结构服务
func NewRedisObject(options CometDB.Options) (*RedisObject, error) {
	db, err := CometDB.Open(options)
	if err != nil {
		return nil, err
	}
	return &RedisObject{db: db}, nil
}

func (r *RedisObject) Del(key []byte) error {
	return r.db.Delete(key)
}

func (r *RedisObject) Type(key []byte) (redisObjectType, error) {
	encValue, err := r.db.Get(key)
	if err != nil {
		return 0, err
	}

	if len(encValue) == 0 {
		return 0, errors.New("value is null")
	}

	// 第一个字节就是类型
	return encValue[0], nil
}
