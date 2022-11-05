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
	fq.Init(50)
	go fq.Tricker()

	success := 0
	for i := 0; i < 2000; i++ {
		ok := fq.Access()
		if ok {
			success += 1
		}
		time.Sleep(time.Millisecond * 5)

		if i%500 == 0 {
			desc := fq.Describe()
			fmt.Println(desc.ToString())
			fmt.Println("+_+_+_+_+_+_+_+_+_+_+_+_+_+_+")
		}
	}
}
