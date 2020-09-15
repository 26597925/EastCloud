package main

import (
	"fmt"
	"sapi/pkg/util/timer"
	"sync/atomic"
	"time"
)

var sum int32 = 0
var N int32 = 300
var mTimer *timer.TimingWheel

func now (param map[string]interface{}) error {
	fmt.Println(time.Now().Format(param["time"].(string)))
	atomic.AddInt32(&sum, 1)
	v := atomic.LoadInt32(&sum)
	if v == 2*N {
		mTimer.Stop()
	}

	return nil
}

func main() {
	mTimer = timer.NewTimingWheel(time.Millisecond * 10)
	mTimer.NewWheel(map[string]interface{}{"time":"2006-01-02 15:04:05"}, time.Millisecond*time.Duration(10), now)
	mTimer.Start()

	for {
		select {

		}
	}

	//var i int32
	//for i = 0; i < N; i++ {
	//	timer.NewWheel(time.Millisecond*time.Duration(10*i), now)
	//	timer.NewWheel(time.Millisecond*time.Duration(10*i), now)
	//}
	//timer.Start()
	//if sum != 2*N {
	//	logger.Info("fail")
	//}
}

