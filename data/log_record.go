package data

type LogRecordPos struct {
	// 文件id，表示将数据存储在哪个文件中
	Fid uint32
	// 偏移，表示将数据存储到了数据文件的哪个位置
	Offset int64
}
