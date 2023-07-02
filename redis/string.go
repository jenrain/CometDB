package redis

import (
	"encoding/binary"
	"time"
)

func (r *RedisObject) Set(key []byte, ttl time.Duration, value []byte) error {
	if value == nil {
		return nil
	}

	// 编码 value : type + expire + payload
	buf := make([]byte, binary.MaxVarintLen64+1)
	buf[0] = String
	var index = 1
	var expire int64 = 0
	if ttl != 0 {
		expire = time.Now().Add(ttl).UnixNano()
	}
	index += binary.PutVarint(buf[index:], expire)

	encValue := make([]byte, index+len(value))
	copy(encValue[:index], buf)
	copy(encValue[index:], value)

	// 调用存储接口写入数据
	return r.db.Put(key, encValue)
}

func (r *RedisObject) Get(key []byte) ([]byte, error) {
	encValue, err := r.db.Get(key)
	if err != nil {
		return nil, err
	}

	// 解码
	dataType := encValue[0]
	if dataType != String {
		return nil, ErrWrongTypeOperation
	}

	var index = 1
	expire, n := binary.Varint(encValue[index:])
	index += n
	// 判断是否过期
	if expire > 0 && expire <= time.Now().UnixNano() {
		return nil, nil
	}
	return encValue[index:], nil
}
