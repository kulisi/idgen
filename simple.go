package idgen

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Simple struct {
	Epoch   int64  // 基础时间
	Machine uint16 // 机器码
	//workerIdBitLength byte   // 机器码长度
	//SeqBitLength      byte   // 序列号长度
	MaxSequence uint32 // 最大序列
	MinSequence uint32 // 最小系列

	TimestampShift byte // 时间戳左移位数
	MachineShift   byte // 机器码左移位数

	Timestamp int64  // 当前序时间（毫秒）
	Sequence  uint32 // 当前序列数

	sync.Mutex // 并发锁(每次获取ID都需要上锁)
}

func NewSimple(opts *Options) *Simple {
	return &Simple{
		Epoch:          opts.BaseTime,
		TimestampShift: (byte)(opts.WorkerIdBitLength + opts.SeqBitLength),

		Timestamp: 0,
		Sequence:  opts.MinSeqNumber,
	}
}

// Next 获取下一个ID
func (s *Simple) Next() int64 {
	// 先上锁，再获取ID
	s.Lock()
	defer s.Unlock()
	// 取当前时间戳
	ct := s.CurrentTimestamp()
	// 比较当前时间戳与最后一次获取ID的时间戳
	if s.Timestamp == ct {
		// 若相等，Sequence自增+1
		s.Sequence++
		if s.Sequence > s.MaxSequence {
			// 若新 Sequence 超过设定的最大值（MaxSequence），等待下一毫秒，并重置 Sequence 为最小值（MinSequence）。
			s.Sequence = s.MinSequence
			ct = s.NextTimestamp()
		}
	} else {
		// 若不相等，重置Sequence为默认最小值
		s.Sequence = s.MinSequence
	}
	// 时间回拨时不操作并输出文本
	if ct < s.Timestamp {
		fmt.Println("Time error for {0} milliseconds", strconv.FormatInt(s.Timestamp-ct, 10))
	}
	// 设置新的 Timestamp
	s.Timestamp = ct
	// 组合ID
	result := int64(ct<<s.TimestampShift) + int64(s.Machine<<s.MachineShift) + int64(s.Sequence)
	// 返回ID
	return result
}

// CurrentTimestamp 获取从起始时间到现在经过的时间
func (s *Simple) CurrentTimestamp() int64 {
	// 获取当前毫秒时间戳
	var millis = time.Now().UnixMilli()
	// 返回起始时间到现在经过的时间
	return millis - s.Epoch
}

// NextTimestamp 获取下一个时间标记（等待下一秒）
func (s *Simple) NextTimestamp() int64 {
	t := s.CurrentTimestamp()
	// 新的时间标记必须大于_LastTimeTick(最后一次获取的时间标记)
	for t <= s.Timestamp {
		// 等待1毫秒后再次尝试获取
		time.Sleep(time.Duration(1) * time.Millisecond)
		t = s.CurrentTimestamp()
	}
	return t
}
