package idgen

type Shift struct {
	*Simple
	TopOverCostCount uint32 // 最大漂移次数

	TurnBackTimestamp      int64 // 时间回拨时间戳
	TurnBackIndex          byte  // 时间回拨次数
	IsOverCost             bool  // 是否发生时间回拨
	OverCostCountInOneTerm uint32
}

// NewShift 创建
func NewShift(opts *Options) ICore {
	return &Shift{
		Simple:                 NewSimple(opts),
		TopOverCostCount:       opts.TopOverCostCount,
		TurnBackTimestamp:      0,
		TurnBackIndex:          0,
		IsOverCost:             false,
		OverCostCountInOneTerm: 0,
	}
}

// Next 获取新的ID
func (s *Shift) Next() int64 {
	s.Lock()
	defer s.Unlock()
	if s.IsOverCost {
		return s.NextOverCostId()
	} else {
		return s.NextNormalId()
	}
}

func (s *Shift) NextOverCostId() int64 {
	ct := s.CurrentTimestamp()
	if ct > s.Timestamp {
		s.Timestamp = ct
		s.Sequence = s.MinSequence
		s.IsOverCost = false
		s.OverCostCountInOneTerm = 0
		return s.CalcId(s.Timestamp)
	}
	if s.OverCostCountInOneTerm >= s.TopOverCostCount {
		s.Timestamp = s.NextTimestamp()
		s.Sequence = s.MinSequence
		s.IsOverCost = false
		s.OverCostCountInOneTerm = 0
		return s.CalcId(s.Timestamp)
	}
	if s.Sequence > s.MaxSequence {
		s.Timestamp++
		s.Sequence = s.MinSequence
		s.IsOverCost = true
		s.OverCostCountInOneTerm++
		return s.CalcId(s.Timestamp)
	}
	return s.CalcId(s.Timestamp)
}

func (s *Shift) NextNormalId() int64 {
	// 取基础时间到现在的时间差，用于下面比较当前时间差与上次获取ID的时间差
	ct := s.CurrentTimestamp()
	// 若少于上次获取ID的时间差，触发时间回避处理机制
	if ct < s.Timestamp {
		// 若为第一次触发时间回拨则初始化时间回拨参数
		if s.TurnBackTimestamp < 1 {
			s.TurnBackTimestamp = s.Timestamp - 1
			s.TurnBackIndex++
			// 每毫秒序列数的前5位是预留位，0用于手工新值，1-4是时间回拨次序
			// 支持4次回拨次序，避免回拨重叠导致ID重复，可无限次回拨，次序循环使用。
			if s.TurnBackIndex > 4 {
				s.TurnBackIndex = 1
			}
		}
		// 返回新ID
		return s.CalcTurnBackId(s.TurnBackTimestamp)
	}
	// 当发生时间回拨并时间追平时，TurnBackTimeTick 清零
	if s.TurnBackTimestamp > 0 {
		s.TurnBackTimestamp = 0
	}
	// 若当前时间差大于上次获取ID的时间差
	if ct > s.Timestamp {
		s.Timestamp = ct
		s.Sequence = s.MinSequence
		// 返回新ID
		return s.CalcId(s.Timestamp)
	}
	// 若新序列数大于允许的最大序列数，则等待下一毫秒。
	if s.Sequence > s.MaxSequence {
		s.Timestamp++
		s.Sequence = s.MinSequence
		s.IsOverCost = true
		s.OverCostCountInOneTerm = 1
		// 返回新ID
		return s.CalcId(s.Timestamp)
	}
	// 返回新ID
	return s.CalcId(s.Timestamp)
}

// CalcId 计算最新ID
func (s *Shift) CalcId(useTimeTick int64) int64 {
	result := int64(useTimeTick<<s.TimestampShift) + int64(s.Machine<<s.MachineShift) + int64(s.Sequence)
	s.Sequence++
	return result
}

// CalcTurnBackId 计算回拨后的最新ID
func (s *Shift) CalcTurnBackId(useTimeTick int64) int64 {
	result := int64(useTimeTick<<s.TimestampShift) + int64(s.Machine<<s.MachineShift) + int64(s.TurnBackIndex)
	s.TurnBackTimestamp--
	return result
}
