package redis

import (
	"CometDB"
	"encoding/binary"
	"errors"
	"math"
	"time"
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

func (r *RedisObject) Close() error {
	return r.db.Close()
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

const (
	maxMetadataSize   = 1 + binary.MaxVarintLen64*2 + binary.MaxVarintLen32
	extraListMetaSize = binary.MaxVarintLen64 * 2

	initialListMark = math.MaxUint64 / 2
)

// 类型元数据
type metaData struct {
	// 数据类型
	dataType byte
	// 过期时间
	expire int64
	//  版本号
	version int64
	// 数据量
	size uint32
	// List 数据结构专用
	head uint64
	tail uint64
}

func (m *metaData) encode() []byte {
	var size = maxMetadataSize
	if m.dataType == List {
		size += extraListMetaSize
	}
	buf := make([]byte, size)

	buf[0] = m.dataType
	index := 1
	index += binary.PutVarint(buf[index:], m.expire)
	index += binary.PutVarint(buf[index:], m.version)
	index += binary.PutVarint(buf[index:], int64(m.size))

	if m.dataType == List {
		index += binary.PutUvarint(buf[index:], m.head)
		index += binary.PutUvarint(buf[index:], m.tail)
	}

	return buf[index:]
}

func decodeMetadata(buf []byte) *metaData {
	dataType := buf[0]

	index := 1
	expire, n := binary.Varint(buf[index:])
	index += n
	version, n := binary.Varint(buf[index:])
	index += n
	size, n := binary.Varint(buf[index:])
	index += n

	var head, tail uint64
	if dataType == List {
		head, n = binary.Uvarint(buf[index:])
		index += n
		tail, _ = binary.Uvarint(buf[index:])
	}

	return &metaData{
		dataType: dataType,
		expire:   expire,
		version:  version,
		size:     uint32(size),
		head:     head,
		tail:     tail,
	}
}

// 查找 key 的类型元数据
func (r *RedisObject) findMetadata(key []byte, dataType redisObjectType) (*metaData, error) {
	metaBuf, err := r.db.Get(key)
	if err != nil && err != CometDB.ErrKeyNotFound {
		return nil, err
	}

	var meta *metaData
	var exist = true
	if err == CometDB.ErrKeyNotFound {
		exist = false
	} else {
		meta = decodeMetadata(metaBuf)
		// 判断数据类型
		if meta.dataType != dataType {
			return nil, ErrWrongTypeOperation
		}
		// 判断过期时间
		if meta.expire != 0 && meta.expire <= time.Now().UnixNano() {
			exist = false
		}
	}

	if !exist {
		meta = &metaData{
			dataType: dataType,
			expire:   0,
			version:  time.Now().UnixNano(),
			size:     0,
		}
		if dataType == List {
			meta.head = initialListMark
			meta.tail = initialListMark
		}
	}
	return meta, nil
}
