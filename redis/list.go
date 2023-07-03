package redis

import (
	"CometDB"
	"encoding/binary"
)

type listInternalKey struct {
	key     []byte
	version int64
	index   uint64
}

func (l *listInternalKey) encode() []byte {
	buf := make([]byte, len(l.key)+8+8)

	// key
	var index = 0
	copy(buf[index:index+len(l.key)], l.key)
	index += len(l.key)

	// version
	binary.LittleEndian.PutUint64(buf[index:index+8], uint64(l.version))
	index += 8

	// index
	binary.LittleEndian.PutUint64(buf[index:], l.index)

	return buf
}

func (r *RedisObject) LPush(key, element []byte) (uint32, error) {
	return r.pushInner(key, element, true)
}

func (r *RedisObject) RPush(key, element []byte) (uint32, error) {
	return r.pushInner(key, element, false)
}

func (r *RedisObject) LPop(key []byte) ([]byte, error) {
	return r.popInner(key, true)
}

func (r *RedisObject) RPop(key []byte) ([]byte, error) {
	return r.popInner(key, false)
}

func (r *RedisObject) LLen(key []byte) (uint32, error) {
	meta, err := r.findMetadata(key, List)
	if err != nil {
		return 0, err
	}
	return meta.size, nil
}

func (r *RedisObject) pushInner(key, element []byte, isLeft bool) (uint32, error) {
	// 查找元数据
	meta, err := r.findMetadata(key, List)
	if err != nil {
		return 0, err
	}

	// 构造数据部分的 key
	l := &listInternalKey{
		key:     key,
		version: meta.version,
	}
	if isLeft {
		l.index = meta.head - 1
	} else {
		l.index = meta.tail
	}

	// 更新元数据和数据部分
	wb := r.db.NewWriteBatch(CometDB.DefaultWriteBatchOptions)
	meta.size++
	if isLeft {
		meta.head--
	} else {
		meta.tail++
	}
	_ = wb.Put(key, meta.encode())
	_ = wb.Put(l.encode(), element)
	if err = wb.Commit(); err != nil {
		return 0, err
	}

	return meta.size, nil
}

func (r *RedisObject) popInner(key []byte, isLeft bool) ([]byte, error) {
	// 查找元数据
	meta, err := r.findMetadata(key, List)
	if err != nil {
		return nil, err
	}
	if meta.size == 0 {
		return nil, nil
	}

	// 构造数据部分的 key
	l := &listInternalKey{
		key:     key,
		version: meta.version,
	}
	if isLeft {
		l.index = meta.head
	} else {
		l.index = meta.tail - 1
	}

	element, err := r.db.Get(l.encode())
	if err != nil {
		return nil, err
	}

	// 更新元数据
	meta.size--
	if isLeft {
		meta.head++
	} else {
		meta.tail--
	}
	if err = r.db.Put(key, meta.encode()); err != nil {
		return nil, err
	}

	return element, nil
}
