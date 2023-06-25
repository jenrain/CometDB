package data

import (
	"encoding/binary"
	"hash/crc32"
)

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// crc type keySize valueSize
// 4 + 1 + 5 + 5
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

// LogRecord 写入到数据文件的记录
// 之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// logRecord的头部信息
type logRecordHeader struct {
	// crc 校验值
	crc uint32

	// 标识 LogRecord 的类型
	recordType LogRecordType

	// key 的长度
	keySize uint32

	// value的长度
	valueSize uint32
}

type LogRecordPos struct {
	// 文件id，表示将数据存储在哪个文件中
	Fid uint32
	// 偏移，表示将数据存储到了数据文件的哪个位置
	Offset int64
}

// EncodeLogRecord 对LogRecord进行编码，返回字节数组及长度
// crc校验 + type类型 + key size + value size + key + value
//
//	4字节     1字节     变长(最大5)   变长(最大5)  变长    变长
func EncodeLogRecord(record *LogRecord) ([]byte, int64) {
	// 初始化一个 header 部分的字节数组
	header := make([]byte, maxLogRecordHeaderSize)

	// 第五个字节存储 Type
	header[4] = record.Type
	var index = 5
	// 5字节之后，存储的是 key 和 value 的长度信息
	// 使用变长类型，节省空间
	index += binary.PutVarint(header[index:], int64(len(record.Key)))
	index += binary.PutVarint(header[index:], int64(len(record.Value)))

	var size = index + len(record.Key) + len(record.Value)
	encBytes := make([]byte, size)

	// 将header部分的内容拷贝过来
	copy(encBytes[:index], header[:index])
	// 将key和value数据拷贝到字节数组中
	copy(encBytes[index:], record.Key)
	copy(encBytes[index+len(record.Key):], record.Value)

	// 对整个 LogRecord 的数据进行crc校验
	crc := crc32.ChecksumIEEE(encBytes[4:])
	// 将crc以小端字节序写入encBytes数组的前四个字节
	binary.LittleEndian.PutUint32(encBytes[:4], crc)

	return encBytes, int64(size)
}

// 对字节数组中的 Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*logRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &logRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]),
		recordType: buf[4],
	}

	var index = 5
	// 取出实际的 key size
	keySize, n := binary.Varint(buf[index:])
	header.keySize = uint32(keySize)
	index += n

	// 取出实际的 value size
	valueSize, n := binary.Varint(buf[index:])
	header.valueSize = uint32(valueSize)
	index += n

	return header, int64(index)
}

// 计算LogRecord的CRC
func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	if lr == nil {
		return 0
	}

	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, lr.Key)
	crc = crc32.Update(crc, crc32.IEEETable, lr.Value)

	return crc
}
