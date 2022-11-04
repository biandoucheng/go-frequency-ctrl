# go-frequency-ctrl

## 概述
```
频率控制

将频率值 qps 分配置到最小10ms的时间段中，当请求进来时对每个时间段的数值进行检查；
当请求数超过当前时间单位ms分配的数值时，向其左右两边的时间段借用，当借用的数值也不够用时该次请求被丢弃；
主频率控制器，会对每个成功访问的请求进行统计，当统计数大于qps时将主控制器置为不可访问；
主控制器的统计周期为1秒，1秒后重置其状态为可访问，统计数重置为0；
```

# 使用方法
[示例代码 <https://github.com/biandoucheng/open-example/tree/main/go-frequency-ctrl>](https://github.com/biandoucheng/open-example/tree/main/go-frequency-ctrl-example)

```
fq := freq.Frequency{}
	fq.Init(40)
	go fq.Tricker()

	success := 0
	for i := 0; i < 50; i++ {
		ok := fq.Access()
		if ok {
			success += 1
		}
		fmt.Println("访问接受 >>", i, ok)
		time.Sleep(time.Millisecond * 2)
	}

	fmt.Printf("%+v", fq)
	fmt.Println("")
	fmt.Println("请求成功数 >>", success)


// 输出
访问接受 >> 0 true
访问接受 >> 1 true
访问接受 >> 2 false
访问接受 >> 3 false
访问接受 >> 4 false
访问接受 >> 5 false
访问接受 >> 6 false
访问接受 >> 7 false
访问接受 >> 8 false
访问接受 >> 9 false
访问接受 >> 10 false
访问接受 >> 11 false
访问接受 >> 12 true
访问接受 >> 13 false
访问接受 >> 14 false
访问接受 >> 15 false
访问接受 >> 16 false
访问接受 >> 17 false
访问接受 >> 18 false
访问接受 >> 19 false
访问接受 >> 20 false
访问接受 >> 21 false
访问接受 >> 22 false
访问接受 >> 23 true
访问接受 >> 24 false
访问接受 >> 25 false
访问接受 >> 26 false
访问接受 >> 27 false
访问接受 >> 28 false
访问接受 >> 29 false
访问接受 >> 30 false
访问接受 >> 31 false
访问接受 >> 32 false
访问接受 >> 33 false
访问接受 >> 34 true
访问接受 >> 35 false
访问接受 >> 36 false
访问接受 >> 37 false
访问接受 >> 38 false
访问接受 >> 39 false
访问接受 >> 40 false
访问接受 >> 41 false
访问接受 >> 42 false
访问接受 >> 43 false
访问接受 >> 44 false
访问接受 >> 45 true
访问接受 >> 46 false
访问接受 >> 47 false
访问接受 >> 48 false
访问接受 >> 49 false
{RWMutex:{w:{state:0 sema:0} writerSem:0 readerSem:0 readerCount:0 readerWait:0} hertz:40 status:0 count:6 busyAfter:32 waves:[{value:1 count:1 status:3} {value:1 count:1 status:3} {value:1 count:1 status:3} {value:1 count:1 status:3} {value:1 count:1 status:3} {value:1 count:1 status:3} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value
:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1} {value:1 count:0 status:1}] waveIndex:4 waveSize:40 wavePartTime:25000}
```