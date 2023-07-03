package redis

import (
	"CometDB"
	"encoding/binary"
)

type hashInternalKey struct {
	key     []byte
	version int64
	field   []byte
}

func (hk *hashInternalKey) encode() []byte {
	buf := make([]byte, len(hk.key)+len(hk.field)+8)
	// key
	var index = 0
	copy(buf[index:index+len(hk.key)], hk.key)
	index += len(hk.key)

	// version
	binary.LittleEndian.PutUint64(buf[index:index+8], uint64(hk.version))
	index += 8

	// field
	copy(buf[index:], hk.field)

	return buf
}

func (r *RedisObject) HSet(key, field, value []byte) (bool, error) {
	// 先查找元数据
	meta, err := r.findMetadata(key, Hash)
	if err != nil {
		return false, err
	}

	// 构造 Hash 数据部分的 key
	hk := &hashInternalKey{
		key:     key,
		version: meta.version,
		field:   field,
	}
	encKey := hk.encode()

	// 先查找是否存在
	var exist = true
	if _, err = r.db.Get(encKey); err == CometDB.ErrKeyNotFound {
		exist = false
	}

	wb := r.db.NewWriteBatch(CometDB.DefaultWriteBatchOptions)
	// 不存在则更新元数据
	if !exist {
		meta.size++
		_ = wb.Put(key, meta.encode())
	}
	_ = wb.Put(encKey, value)
	if err = wb.Commit(); err != nil {
		return false, err
	}
	return !exist, nil
}

func (r *RedisObject) HGet(key, field []byte) ([]byte, error) {
	meta, err := r.findMetadata(key, Hash)
	if err != nil {
		return nil, err
	}
	if meta.size == 0 {
		return nil, nil
	}

	hk := &hashInternalKey{
		key:     key,
		version: meta.version,
		field:   field,
	}

	return r.db.Get(hk.encode())
}

func (r *RedisObject) HDel(key, field []byte) (bool, error) {
	meta, err := r.findMetadata(key, Hash)
	if err != nil {
		return false, err
	}
	if meta.size == 0 {
		return false, nil
	}

	hk := &hashInternalKey{
		key:     key,
		version: meta.version,
		field:   field,
	}
	encKey := hk.encode()

	// 先查看是否存在
	var exist = true
	if _, err = r.db.Get(encKey); err == CometDB.ErrKeyNotFound {
		exist = false
	}

	if exist {
		wb := r.db.NewWriteBatch(CometDB.DefaultWriteBatchOptions)
		meta.size--
		_ = wb.Put(key, meta.encode())
		_ = wb.Delete(encKey)
		if err = wb.Commit(); err != nil {
			return false, err
		}
	}

	return exist, nil
}

func (r *RedisObject) HLen(key []byte) (uint32, error) {
	meta, err := r.findMetadata(key, Hash)
	if err != nil {
		return 0, err
	}
	return meta.size, nil
}
