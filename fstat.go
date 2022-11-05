package frequency

import (
	"sync"
	"time"
)

// 访问统计
// 访问频率统计,将15分钟按5秒钟一段分180份
type FrequencyStat struct {
	sync.RWMutex // 读写锁

	SecondAccessStats   [2]int   // 访问数
	SecondPassedStats   [2]int   // 通过数
	FTSecondAccessStats [181]int // 访问数
	FTSecondPassedStats [181]int // 通过数

	FStatus  int // 频控器状态 秒级
	FWStatus int // 频控波段状态  (1~10)ms 级
}

// 做索引重置
func (f *FrequencyStat) Tricker() {
	tricker := time.NewTicker(time.Second * 1)
	second := 0

	for {
		<-tricker.C
		second += 1

		f.Lock()

		f.SecondAccessStats[1] = f.SecondAccessStats[0]
		f.SecondAccessStats[0] = 0
		f.SecondPassedStats[1] = f.SecondPassedStats[0]
		f.SecondPassedStats[0] = 0

		if second%5 == 0 {
			for i := 179; i >= 0; i-- {
				if i > 0 {
					f.FTSecondAccessStats[i] = f.FTSecondAccessStats[i-1]
					f.FTSecondPassedStats[i] = f.FTSecondPassedStats[i-1]
				} else {
					f.FTSecondAccessStats[i] = 0
					f.FTSecondPassedStats[i] = 0
				}
			}
		}

		f.Unlock()
	}
}

func (f *FrequencyStat) Describe() FrequencyDescribe {
	defer f.RUnlock()
	f.RLock()

	desc := FrequencyDescribe{}

	// 1 秒钟统计
	desc.SecnodAccess = f.SecondAccessStats[1]
	desc.SecondPassed = f.SecondPassedStats[1]

	// 1 分钟统计
	for i := 1; i < 13; i++ {
		desc.MinuteAccess += f.FTSecondAccessStats[i]
		desc.MinutePassed += f.FTSecondPassedStats[i]
	}

	// 5 分钟统计
	for i := 1; i < 61; i++ {
		desc.FMinuteAccess += f.FTSecondAccessStats[i]
		desc.FMinutePassed += f.FTSecondPassedStats[i]
	}

	// 15 分钟统计
	for i := 1; i < 181; i++ {
		desc.FTMinuteAccess += f.FTSecondAccessStats[i]
		desc.FTMinutePassed += f.FTSecondPassedStats[i]
	}

	// 计算通过率
	desc.CalculateRate()

	// 记录状态
	desc.SetStatus(f.FStatus, f.FWStatus)

	return desc
}

// 统计
func (f *FrequencyStat) Stat(pass bool, fst, fwst int) {
	f.Lock()

	f.SecondAccessStats[0] += 1
	f.FTSecondAccessStats[0] += 1

	if pass {
		f.SecondPassedStats[0] += 1
		f.FTSecondPassedStats[0] += 1
	}

	f.FStatus = fst
	f.FWStatus = fwst

	f.Unlock()
}
