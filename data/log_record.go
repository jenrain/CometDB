package data

import "encoding/binary"

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
func EncodeLogRecord(record *LogRecord) ([]byte, int64) {
	return nil, 0
}

// 对字节数组中的 Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*logRecordHeader, int64) {
	return nil, 0
}

// 计算LogRecord的CRC
func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	return 0
}
