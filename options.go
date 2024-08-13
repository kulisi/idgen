package idgen

type Options struct {
	Method            uint16 // 雪花计算方法（0.传统算法|1.漂移算法）
	BaseTime          int64  // 基础时间（毫秒单位），要求：不能超过当前系统时间
	WorkerId          uint16 // 机器码，必须由外部设定，最大值 2^WorkerIdBitLength-1
	WorkerIdBitLength byte   // 机器码位长，取值范围[1,15]，要求：序列号长度+机器码长度<=22
	SeqBitLength      byte   // 序列数位长，取值范围[3,21]，要求：序列号长度+机器码长度<=22
	MaxSeqNumber      uint32 // 最大序列数，取值范围[MinSeqNumber,2^SeqBitLength-1]
	MinSeqNumber      uint32 // 最小序列数，取值范围[5,MaxSeqNumber],
	TopOverCostCount  uint32 // 最大漂移次数，取值范围[500,10000]
}

func DefaultOptions(machine uint16) *Options {
	return &Options{
		Method:            1,             // 默认值 漂移算法
		WorkerId:          machine,       // 函数参数传入值
		BaseTime:          1723132800000, // 开始时间 2024-08-09
		WorkerIdBitLength: 6,             // 默认值
		SeqBitLength:      6,             // 默认值
		MaxSeqNumber:      (1 << 6) - 1,  // 默认值
		MinSeqNumber:      5,             // 每毫秒的前5个序列数对应0-4是保留位，其中1-4是时间回拨相应预留位，0是手工新值预留位
		TopOverCostCount:  2000,          // 默认值
	}
}
