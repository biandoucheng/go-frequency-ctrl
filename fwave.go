package frequency

// 波段
type WaveBand struct {
	value  int // 波段值
	count  int // 计数
	status int // 状态 0 无效的，1 正常的，2 忙碌的，3 禁止访问
}

// 初始化一个波段
func (w *WaveBand) Init(val int) {
	w.value = val
	w.status = vStatusNormal
}

// 重置波段信息
func (w *WaveBand) ReSet() {
	w.status = vStatusNormal
	w.count = 0
}

// 是否正常
func (w *WaveBand) IsNormal() bool {
	return IsVNormal(w.status)
}

// 是否忙碌的
func (w *WaveBand) IsBusy() bool {
	return IsVBusy(w.status)
}

// 是否禁止访问的
func (w *WaveBand) IsForbidden() bool {
	return IsVForbidden(w.status)
}

// 判断可访问性没，并递减值
func (w *WaveBand) Access() bool {
	if w.IsForbidden() {
		return false
	}

	if w.value <= w.count {
		w.status = vStatusForbidden
		return false
	}

	w.count += 1
	return true
}
