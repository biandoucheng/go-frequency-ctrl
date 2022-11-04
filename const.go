package frequency

// 频率状态 0 正常，1 忙碌的，2 禁止访问
const (
	// 正常状态
	fStatusNormal = 0
	// 忙碌
	fStatusBusy = 1
	// 禁止访问
	fStatusForbidden = 2
)

// 是否是正常状态
func IsFNormal(st int) bool {
	return st == fStatusNormal
}

// 是否禁止访问
func IsFForbidden(st int) bool {
	return st == fStatusForbidden
}

// 忙碌的
func IsFBusy(st int) bool {
	return st == fStatusBusy
}

// 波段状态 1 正常的，2 忙碌的，3 禁止访问
const (
	// 无效状态
	vStatusInvalid = 0
	// 正常状态
	vStatusNormal = 1
	// 忙碌的
	vStatusBusy = 2
	// 禁止访问
	vStatusForbidden = 3
)

// 是否无效的
func IsVInvalid(st int) bool {
	return st == vStatusInvalid
}

// 是否正常的
func IsVNormal(st int) bool {
	return st == vStatusNormal
}

// 是否忙碌的
func IsVBusy(st int) bool {
	return st == vStatusBusy
}

// 是否禁止的
func IsVForbidden(st int) bool {
	return st == vStatusForbidden
}
