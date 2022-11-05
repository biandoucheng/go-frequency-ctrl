# go-frequency-ctrl

## 概述
```
频率控制

将频率值 qps 分配置到最小1ms的时间段中，当请求进来时对每个时间段的数值进行检查；
当请求数超过当前时间单位ms分配的数值时，向其左右两边的时间段借用，当借用的数值也不够用时该次请求被丢弃；
主频率控制器，会对每个成功访问的请求进行统计，当统计数大于qps时将主控制器置为不可访问；
主控制器的统计周期为1秒，1秒后重置其状态为可访问，统计数重置为0；
```

## 结构体及方法
```
// 频率
type Frequency struct {
	sync.RWMutex // 读写锁

	hertz  int // 赫兹，频率，每秒次数
	status int // 状态 0 正常，1 忙碌的，2 禁止访问
	count  int // 一秒内的请求数

	busyAfter int // 忙碌阈值

	waves        []WaveBand // 波段
	waveIndex    int        // 波段索引
	waveSize     int        // 波长
	wavePartTime int        // 波的单位时长 mrocs 微妙

	stat FrequencyStat // 统计
}

该结构体是频率控制器的主结构体：
Frequency.hertz 设置基本频率值（times/second），Init 方法中会将该频率分配到每个单位毫秒时间段中去。
实际控制访问通过的与否是根据单位毫秒是否能够争取到足够的访问允许次数，具体见 Frequency.Access 方法介绍。

Frequency.status 是主频控器的状态值，当 count 值大于等于 hertz 时该状态会被置为禁止访问，此时访问会被直接拒绝，否则会去询问当前毫秒单位下是否有足够的访问次数有则通过。
其中 忙碌 状态不会影响访问，他只是给使用者观察频控压力的一个指标。

Frequency.count 是一秒中的通过的访问次数，每经过1秒循环会被清零一次（清零的同时会把status置为正常状态 。

Frequency.busyAfter 它是用来判断频控器是否处于忙碌状态，等于 hertz * 0.8 取整，当 count >= busyAfter 时 status 会被置为忙碌状态。

Frequency.waves 波段，它是按照一定规则将 hertz 的值分配到单位(1~10)ms下的一组波段。每个波段有记录它允许访问的次数以及允许访问的状态。

Frequency.waveIndex 波段索引，它指明当前访问与否应该向哪个波段询问，每经历一个单位 ms 循环会自动加一，当完成1秒循环会被重置为0。

Frequency.waveSize 波段数量，它是在1秒内按照单位 ms 将hertz分割的波段数量的记录，根据他在循环中判断是否完成一整个循环。

Frequency.wavePartTime 单位毫秒（>=1ms < 1s），用来判断多久需要移动一次波段索引，实际上它存的是微妙值，为的是最大程度避免波段分配时毫秒除不尽带来的时间循环不准确问题，它时间已经控制到不到 0.1 毫秒的偏差之内了。

Frequency.stat 频控统计对象，他用来统计访问频率和成功率以及当前频控器状态。

Frequency.Init(hertz int) 初始化一个频控器。

Frequency.Access() 访问允许判断，它会先判断当前波段是否允许访问，如果不允许则会分别向左右波段（如果有）借访问次数（有借无还），相邻的这三个波段只要任意一个允许访问就会返回允许访问，否则不允许访问。

Frequency.Tricker() 定时器，完成一系列状态重置以及波段索引移动任务。

Frequency.Describe() 输出频控器的当前统计信息。


;;
// 波段
type WaveBand struct {
	value  int // 波段值
	count  int // 计数
	status int // 状态 0 无效的，1 正常的，2 忙碌的，3 禁止访问
}

WaveBand.value 当前波段允许的访问次数。

WaveBand.count 当前波段已通过访问次数。

WaveBand.status 当前波段状态。

WaveBand.Init(val int) 初始化一个波段，val 是该波段允许的访问次数。

WaveBand.Access() 当前波段允许访问判断。


;;
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

FrequencyStat.SecondAccessStats 当前秒以及前一秒的访问数。

FrequencyStat.SecondPassedStats 当前秒以及前一秒的访问通过数。

FrequencyStat.FTSecondAccessStats 最近15分钟的每5秒访问数，它实际统计了当前5秒以及过去15分钟的每5秒。

FrequencyStat.FTSecondPassedStats 最近15分钟的每5秒访问通过数，它实际统计了当前5秒以及过去15分钟的每5秒。

FrequencyStat.FStatus 当前主频控器的状态。

FrequencyStat.FWStatus 当前波段状态。

FrequencyStat.Stat(pass bool, fst, fwst int) 访问统计，pass 访问是否通过，fst 主频控器状态，fwst 当前波段状态。

Frequency.Tricker() 定时器，负责对统计数据进行位移操作，实现每次只对索引0的值进行累加，定时对数组内容整体右移一位并对索引0的值置0。

FrequencyStat.Describe() 统计信息输出，秒级统计输出前一秒，分钟级别统计为5秒前的数据。


;;
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

FrequencyDescribe.CalculateRate() 计算通过率。

FrequencyDescribe.SetStatus(fst, fwst int) 设置频控器状态，fst 主频控器状态，fwst 当前波段状态。

FrequencyDescribe.ToString() 以字符串形式输出统计信息。
```

## 使用方法
[示例代码 <https://github.com/biandoucheng/open-example/tree/main/go-frequency-ctrl-example>](https://github.com/biandoucheng/open-example/tree/main/go-frequency-ctrl-example)

```
import (
	"fmt"
	"testing"
	"time"

	freq "github.com/biandoucheng/go-frequency-ctrl"
)

func TestFrequencyAccess(t *testing.T) {
	fq := freq.Frequency{}
	fq.Init(50)
	go fq.Tricker()

	success := 0
	for i := 0; i < 1000; i++ {
		ok := fq.Access()
		if ok {
			success += 1
		}
		time.Sleep(time.Millisecond * 5)

		if i%1000 == 0 {
			desc := fq.Describe()
			fmt.Println(desc.ToString())
		}
	}
}


// 输出
=== RUN   TestFrequencyAccess

        QP/S:
        SecnodAccess: 0
        SecondPassed: 0
        SecondPassRate: 0.00

        QP/M
        MinuteAccess: 0
        MinutePassed: 0
        MinutePassRate: 0.00

        QP/5M:
        FMinuteAccess: 0
        FMinutePassed: 0
        FMinutePassRate: 0.00

        QP/15M:
        FTMinuteAccess: 0
        FTMinutePassed: 0
        FTMinutePassRate: 0.00

        FStatus:
        FStatus: NORMAL
        FWStatus: NORMAL

+_+_+_+_+_+_+_+_+_+_+_+_+_+_+

        QP/S:
        SecnodAccess: 179
        SecondPassed: 50
        SecondPassRate: 27.93

        QP/M
        MinuteAccess: 0
        MinutePassed: 0
        MinutePassRate: 0.00

        QP/5M:
        FMinuteAccess: 0
        FMinutePassed: 0
        FMinutePassRate: 0.00

        QP/15M:
        FTMinuteAccess: 0
        FTMinutePassed: 0
        FTMinutePassRate: 0.00

        FStatus:
        FStatus: NORMAL
        FWStatus: FORBIDDEN

+_+_+_+_+_+_+_+_+_+_+_+_+_+_+

        QP/S:
        SecnodAccess: 179
        SecondPassed: 50
        SecondPassRate: 27.93

        QP/M
        MinuteAccess: 896
        MinutePassed: 250
        MinutePassRate: 27.90

        QP/5M:
        FMinuteAccess: 896
        FMinutePassed: 250
        FMinutePassRate: 27.90

        QP/15M:
        FTMinuteAccess: 896
        FTMinutePassed: 250
        FTMinutePassRate: 27.90

        FStatus:
        FStatus: NORMAL
        FWStatus: FORBIDDEN

+_+_+_+_+_+_+_+_+_+_+_+_+_+_+

        QP/S:
        SecnodAccess: 180
        SecondPassed: 50
        SecondPassRate: 27.78

        QP/M
        MinuteAccess: 896
        MinutePassed: 250
        MinutePassRate: 27.90

        QP/5M:
        FMinuteAccess: 896
        FMinutePassed: 250
        FMinutePassRate: 27.90

        QP/15M:
        FTMinuteAccess: 896
        FTMinutePassed: 250
        FTMinutePassRate: 27.90

        FStatus:
        FStatus: NORMAL
        FWStatus: FORBIDDEN

+_+_+_+_+_+_+_+_+_+_+_+_+_+_+
--- PASS: TestFrequencyAccess (11.15s)
```