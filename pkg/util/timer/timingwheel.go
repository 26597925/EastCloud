package timer

import (
	"container/list"
	"sync"
	"time"
)

//referer https://github.com/cloudwu/skynet/blob/master/skynet-src/skynet_timer.c

const (
	TimeNearShift  = 8
	TimeNear       = 1 << TimeNearShift
	TimeLevelShift = 6
	TimeLevel      = 1 << TimeLevelShift
	TimeNearMask   = TimeNear - 1
	TimeLevelMask  = TimeLevel - 1
)

type TimingWheel struct {
	near [TimeNear]*list.List
	t    [4][TimeLevel]*list.List
	sync.Mutex
	time uint32
	tick time.Duration
	quit chan bool
}

type ExecuteHandler func(param map[string]interface{}) error

type Wheel struct {
	param   map[string]interface{}
	expire  uint32
	handler ExecuteHandler
}

func NewTimingWheel(d time.Duration) *TimingWheel {
	quit := make(chan bool, 1)

	tw := &TimingWheel{
		time: 0,
		tick: d,
		quit: quit,
	}

	var i, j int
	for i = 0; i < TimeNear; i++ {
		tw.near[i] = list.New()
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < TimeLevel; j++ {
			tw.t[i][j] = list.New()
		}
	}

	return tw
}

func (tw *TimingWheel) Start() {
	go func() {
		tick := time.NewTicker(tw.tick)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				tw.update()
			case <-tw.quit:
				return
			}
		}
	}()
}

func (tw *TimingWheel) Stop() {
	close(tw.quit)
}

func (tw *TimingWheel) NewWheel (p map[string]interface{}, d time.Duration, f ExecuteHandler) *Wheel {
	w := &Wheel{
		param: p,
		expire: uint32(d/tw.tick) + tw.time,
		handler: f,
	}

	tw.Lock()
	tw.addWheel(w)
	tw.Unlock()
	return w
}

func (tw *TimingWheel) addWheel(w *Wheel) {
	if (w.expire | TimeNearMask) == (tw.time | TimeNearMask) {
		tw.near[w.expire & TimeNearMask].PushBack(w)
	} else {
		var i uint32
		var mask uint32 = TimeNear << TimeLevelShift
		for i = 0; i < 3; i++ {
			if (w.expire | (mask - 1)) == (tw.time | (mask - 1)) {
				break
			}
			mask <<= TimeLevelShift
		}

		tw.t[i][(w.expire >> (TimeNearShift + i*TimeLevelShift)) & TimeLevelMask].PushBack(w)
	}
}

func (tw *TimingWheel) moveList(level, idx int) {
	vec := tw.t[level][idx]
	front := vec.Front()
	vec.Init()
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Wheel)
		tw.addWheel(node)
	}
}

func (tw *TimingWheel) shift() {
	tw.Lock()
	var mask uint32 = TimeNear
	tw.time++
	ct := tw.time
	if ct == 0 {
		tw.moveList(3, 0)
	} else {
		time := ct >> TimeNearShift
		var i int = 0
		for (ct & (mask - 1)) == 0 {
			idx := int(time & TimeLevelMask)
			if idx != 0 {
				tw.moveList(i, idx)
				break
			}
			mask <<= TimeLevelShift
			time >>= TimeLevelShift
			i++
		}
	}
	tw.Unlock()
}

func (tw *TimingWheel) execute() {
	tw.Lock()
	idx := tw.time & TimeNearMask
	vec := tw.near[idx]
	if vec.Len() > 0 {
		front := vec.Front()
		vec.Init()
		tw.Unlock()
		// dispatch_list don't need lock
		for e := front; e != nil; e = e.Next() {
			node := e.Value.(*Wheel)
			go node.handler(node.param)
		}
		return
	}

	tw.Unlock()
}

func (tw *TimingWheel) update() {
	// try to dispatch timeout 0 (rare condition)
	tw.execute()

	// shift time first, and then dispatch timer message
	tw.shift()

	tw.execute()

}

