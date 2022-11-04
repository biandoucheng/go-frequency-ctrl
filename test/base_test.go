package test

// 基础测试

import (
	"fmt"
	"testing"
	"time"

	freq "github.com/biandoucheng/go-frequency-ctrl"
)

func TestFrequencyInitValuue(t *testing.T) {
	fq := freq.Frequency{}
	fq.Init(300)
	fmt.Printf("%+v", fq)
}

func TestFrequencyAccess(t *testing.T) {
	fq := freq.Frequency{}
	fq.Init(100)
	go fq.Tricker()

	success := 0
	for i := 0; i < 200; i++ {
		ok := fq.Access()
		if ok {
			success += 1
		}
		fmt.Println("访问接受 >>", i, ok)
		time.Sleep(time.Millisecond * 2)
	}

	fmt.Printf("%+v", fq)
	fmt.Println("")
	fmt.Println("当前时间戳 >>", time.Now().Unix())
	fmt.Println("请求成功数 >>", success)
}
