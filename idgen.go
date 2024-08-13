package idgen

import (
	"errors"
	"fmt"
	"time"
)

// MinTime 基础时间最小取值
const MinTime = int64(631123200000)

var (
	ErrOptions           = errors.New("options error")
	ErrBaseTime          = errors.New("base time error (range:['1990-01-01',time.Now()])")
	ErrSeqBitLength      = errors.New("sequence bit length error (range:[2,21])")
	ErrWorkerIdBitLength = errors.New("worker bit length error (range:[1,21])")
	ErrWorkerId          = func(i uint16) error { return fmt.Errorf("worker id error (range:[0, %d])", i) }
	ErrMaxSeqNumber      = func(i uint32) error { return fmt.Errorf("max seq number error (range:[1, %d])", i) }
	ErrMinSeqNumber      = func(i uint32) error { return fmt.Errorf("min seq number error (range:[5, %d])", i) }
	ErrTopOverCostCount  = errors.New("top over cost count error (range:[0, 10000])")
)

type IdGenerator struct {
	core ICore
	opts *Options
}

func NewIdGenerator(options *Options) (*IdGenerator, error) {
	if options == nil {
		return nil, ErrOptions
	}
	// 验证 BaseTime (基础时间不能小于 1990-01-01 也不能大于当前时间)
	if options.BaseTime < MinTime || options.BaseTime > time.Now().UnixMilli() {
		return nil, ErrBaseTime
	}
	// 验证 SeqBitLength (取值范围:[2,21])
	if options.SeqBitLength < 2 || options.SeqBitLength > 21 {
		return nil, ErrSeqBitLength
	}
	// 验证 WorkerIdBitLength (取值范围：[1, 21])
	if options.WorkerIdBitLength < 1 || options.WorkerIdBitLength > 22 {
		return nil, ErrWorkerIdBitLength
	}
	// 验证 WorkerIdBitLength + SeqBitLength (取值范围：[1, 22])
	if options.WorkerIdBitLength+options.SeqBitLength > 22 {
		return nil, ErrWorkerIdBitLength
	}
	// 验证 WorkerId (取值范围：[0, (1<<WorkerIdBitLength)-1]) 不能有默认值
	if options.WorkerId < 0 || options.WorkerId > uint16(1<<options.WorkerIdBitLength)-1 {
		return nil, ErrWorkerId(uint16(1<<options.WorkerIdBitLength) - 1)
	}
	// 验证 MaxSeqNumber
	if options.MaxSeqNumber < 0 || options.MaxSeqNumber > uint32(1<<options.SeqBitLength)-1 {
		return nil, ErrMaxSeqNumber(uint32(1<<options.SeqBitLength) - 1)
	}
	// 验证MinSeqNumber
	if options.MinSeqNumber < 5 || options.MinSeqNumber > uint32(1<<options.SeqBitLength)-1 {
		return nil, ErrMinSeqNumber(uint32(1<<options.SeqBitLength) - 1)
	}
	// 验证 TopOverCostCount
	if options.TopOverCostCount < 0 || options.TopOverCostCount > 10000 {
		return nil, ErrTopOverCostCount
	}

	var core ICore
	switch options.Method {
	case 0:
		core = NewSimple(options)
	case 1:
		core = NewShift(options)
	default:
		core = NewSimple(options)
	}
	return &IdGenerator{
		core: core,
		opts: options,
	}, nil
}

// NewID 生成一个唯一ID
func (n *IdGenerator) NewID() int64 {
	return n.core.Next()
}

// ExtractTime 从ID中提取时间
func (n *IdGenerator) ExtractTime(id int64) time.Time {
	return time.UnixMilli(id>>(n.opts.WorkerIdBitLength+n.opts.SeqBitLength) + n.opts.BaseTime)
}

// ExtractSeq 从ID中提取序列
//func (n *IdGenerator) ExtractSeq(id int64) time.Time {
//	return time.UnixMilli(id>>(n.Opts.WorkerIdBitLength+n.Opts.SeqBitLength) + n.Opts.BaseTime)
//}
