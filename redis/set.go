package redis

import (
	"CometDB"
	"encoding/binary"
)

type setInternalKey struct {
	key     []byte
	version int64
	member  []byte
}

func (s *setInternalKey) encode() []byte {
	buf := make([]byte, len(s.key)+len(s.member)+8+4)
	// key
	var index = 0
	copy(buf[index:index+len(s.key)], s.key)
	index += len(s.key)

	// version
	binary.LittleEndian.PutUint64(buf[index:index+8], uint64(s.version))
	index += 8

	// member
	copy(buf[index:index+len(s.member)], s.member)
	index += len(s.member)

	// member size
	binary.LittleEndian.PutUint32(buf[index:], uint32(len(s.member)))

	return buf
}

func (r *RedisObject) SAdd(key, member []byte) (bool, error) {
	// 查找元数据
	meta, err := r.findMetadata(key, Set)
	if err != nil {
		return false, err
	}

	// 构造一个数据部分的 key
	s := &setInternalKey{
		key:     key,
		version: meta.version,
		member:  member,
	}

	var ok bool
	if _, err = r.db.Get(s.encode()); err == CometDB.ErrKeyNotFound {
		// 不存在的话则更新
		wb := r.db.NewWriteBatch(CometDB.DefaultWriteBatchOptions)
		meta.size++
		_ = wb.Put(key, meta.encode())
		_ = wb.Put(s.encode(), nil)
		if err = wb.Commit(); err != nil {
			return false, err
		}
		ok = true
	}

	return ok, nil
}

func (r *RedisObject) SIsMember(key, member []byte) (bool, error) {
	meta, err := r.findMetadata(key, Set)
	if err != nil {
		return false, err
	}
	if meta.size == 0 {
		return false, nil
	}

	// 构造一个数据部分的 key
	s := &setInternalKey{
		key:     key,
		version: meta.version,
		member:  member,
	}

	_, err = r.db.Get(s.encode())
	if err != nil && err != CometDB.ErrKeyNotFound {
		return false, err
	}
	if err == CometDB.ErrKeyNotFound {
		return false, nil
	}
	return true, nil
}

func (r *RedisObject) SRem(key, member []byte) (bool, error) {
	meta, err := r.findMetadata(key, Set)
	if err != nil {
		return false, err
	}
	if meta.size == 0 {
		return false, nil
	}

	// 构造一个数据部分的 key
	s := &setInternalKey{
		key:     key,
		version: meta.version,
		member:  member,
	}

	if _, err = r.db.Get(s.encode()); err == CometDB.ErrKeyNotFound {
		return false, nil
	}

	// 更新
	wb := r.db.NewWriteBatch(CometDB.DefaultWriteBatchOptions)
	meta.size--
	_ = wb.Put(key, meta.encode())
	_ = wb.Delete(s.encode())
	if err = wb.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisObject) SCard(key []byte) (uint32, error) {
	meta, err := r.findMetadata(key, Set)
	if err != nil {
		return 0, err
	}
	return meta.size, nil
}
