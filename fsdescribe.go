package frequency

import "fmt"

// 统计描述
type FrequencyDescribe struct {
	SecnodAccess   int     `json:"secondAccess"`   // 秒访问数
	SecondPassed   int     `json:"secondPassed"`   // 秒通过量
	SecondPassRate float32 `json:"secondPassRate"` // 秒通过率

	MinuteAccess   int     `json:"minuteAccess"`   // 分钟访问数
	MinutePassed   int     `json:"minutePassed"`   // 分钟通过量
	MinutePassRate float32 `json:"minutePassRate"` // 分钟通过率

	FMinuteAccess   int     `json:"fMinuteAccess"`   // 5分钟访问量
	FMinutePassed   int     `json:"fMinutePassed"`   // 5分钟通过量
	FMinutePassRate float32 `json:"fMinutePassRate"` // 5分钟通过率

	FTMinuteAccess   int     `json:"fTMinuteAccess"`   // 15分钟访问量
	FTMinutePassed   int     `json:"fTMinutePassed"`   // 15分钟通过量
	FTMinutePassRate float32 `json:"fTMinutePassRate"` // 15分钟通过率

	FStatus  string `json:"fStatus"`  // 频控器状态 秒级
	FWStatus string `json:"fWStatus"` // 频控波段状态  (1~10)ms 级
}

// 计算通过
func (f *FrequencyDescribe) CalculateRate() {
	if f.SecondPassed <= 0 {
		f.SecondPassRate = 0
	} else {
		f.SecondPassRate = float32(f.SecondPassed) / float32(f.SecnodAccess) * 100
	}

	if f.MinutePassed <= 0 {
		f.MinutePassRate = 0
	} else {
		f.MinutePassRate = float32(f.MinutePassed) / float32(f.MinuteAccess) * 100
	}

	if f.FMinutePassed <= 0 {
		f.FMinutePassRate = 0
	} else {
		f.FMinutePassRate = float32(f.MinutePassed) / float32(f.MinuteAccess) * 100
	}

	if f.FTMinutePassed <= 0 {
		f.FTMinutePassRate = 0
	} else {
		f.FTMinutePassRate = float32(f.FTMinutePassed) / float32(f.FTMinuteAccess) * 100
	}
}

// 设置状态
func (f *FrequencyDescribe) SetStatus(fst, fwst int) {
	f.FStatus = FStatusToString(fst)
	f.FWStatus = VStatusToString(fwst)
}

//
func (f *FrequencyDescribe) ToString() string {
	val := `
	QP/S:
	SecnodAccess: %d
	SecondPassed: %d
	SecondPassRate: %.2f

	QP/M
	MinuteAccess: %d
	MinutePassed: %d
	MinutePassRate: %.2f

	QP/5M:
	FMinuteAccess: %d
	FMinutePassed: %d
	FMinutePassRate: %.2f

	QP/15M:
	FTMinuteAccess: %d
	FTMinutePassed: %d
	FTMinutePassRate: %.2f

	FStatus:
	FStatus: %s
	FWStatus: %s
	`

	return fmt.Sprintf(
		val,
		f.SecnodAccess, f.SecondPassed, f.SecondPassRate,
		f.MinuteAccess, f.MinutePassed, f.MinutePassRate,
		f.FMinuteAccess, f.FMinutePassed, f.FMinutePassRate,
		f.FTMinuteAccess, f.FTMinutePassed, f.FTMinutePassRate,
		f.FStatus, f.FWStatus)
}
