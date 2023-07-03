package redis

import "errors"

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
