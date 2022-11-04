package frequency

import (
	"sync"
	"time"
)

// 频率
type Frequency struct {
	sync.RWMutex // 写锁

	hertz  int // 赫兹，频率，每秒次数
	status int // 状态 0 正常，1 忙碌的，2 禁止访问
	count  int // 一秒内的请求数

	busyAfter int // 忙碌阈值

	waves        []WaveBand // 波段将频率分割成10个波段 /100ms
	waveIndex    int        // 波段索引
	waveSize     int        // 波长
	wavePartTime int        // 波的单位时长 mrocs 微妙
}

// 初始化频率控制器
func (f *Frequency) Init(hz int) {
	defer f.Unlock()
	f.Lock()

	// 初始化名称频率
	f.hertz = hz

	if f.hertz == 0 {
		f.status = fStatusForbidden
		return
	}

	// 初始化状态
	f.status = fStatusNormal

	// 初始化频率阈值
	f.busyAfter = int(float32(f.hertz) * 0.80)

	// 计算波段
	if f.hertz <= 100 {
		f.waveSize = f.hertz
	} else {
		f.waveSize = 100
	}
	f.wavePartTime = 1000000 / f.waveSize

	// 填充波段
	part := 1
	remain := 0
	if f.hertz > 100 {
		part = f.hertz / 100
		remain = f.hertz % 100
	}

	f.waves = make([]WaveBand, f.waveSize)
	for i := 0; i < f.waveSize; i++ {
		tmp := part
		if remain > 0 {
			tmp += 1
			remain -= 1
		}

		w := WaveBand{}
		w.Init(tmp)
		f.waves[i] = w
	}

}

// 是否禁止访问
func (f *Frequency) IsForbidden() bool {
	return IsFForbidden(f.status)
}

// 设置正常状态
func (f *Frequency) SetNormal() {
	f.status = fStatusNormal
}

// 设置忙碌
func (f *Frequency) SetBusy() {
	f.status = fStatusBusy
}

// 设置禁止访问
func (f *Frequency) SetForbidden() {
	f.status = fStatusForbidden
}

// 访问
func (f *Frequency) Access() bool {
	// 加锁
	defer f.Unlock()
	f.Lock()

	// 访问禁止
	if f.IsForbidden() {
		return false
	}

	// 推波助澜，采用qpms借调方式控制波峰
	idx := f.waveIndex

	// 当前波段判断
	if f.waves[idx].Access() {
		f.count += 1
		return true
	}

	bfidx := idx - 1
	if bfidx > 0 && f.waves[bfidx].Access() {
		f.count += 1
		return true
	}

	afidx := idx + 1
	if afidx+1 <= f.waveSize && f.waves[afidx].Access() {
		f.count += 1
		return true
	}

	return false
}

// 定时器完成一系列重置工作
func (f *Frequency) Tricker() {
	if f.IsForbidden() {
		return
	}

	tricker := time.NewTicker(time.Microsecond * time.Duration(f.wavePartTime))
	cycle_count := 0

	// 循环
	for {
		// 定时
		<-tricker.C

		// 锁定
		f.Lock()

		// 循环计数
		cycle_count += 1

		if cycle_count == f.waveSize {
			// 重置waves
			for i := 0; i < f.waveSize; i++ {
				f.waves[i].ReSet()
			}

			// 重置wave 索引
			f.waveIndex = 0
			cycle_count = 0

			// 重置计数
			f.count = 0
			f.SetNormal()
		} else {
			// 判断频率是否超了
			if f.count > f.hertz {
				f.SetForbidden()
			}

			// 忙碌判断
			if f.count > f.busyAfter {
				f.SetBusy()
			}

			f.waveIndex += 1
		}

		// 解锁
		f.Unlock()
	}
}
