//author gonghh
//qq 29185807
//mail 29185807@qq.com
//Copyright 2014 gonghh. All rights reserved.
package routinebalance

import (
	"sync"
	"sync/atomic"
)

var func_chan chan func()
var routine_count int32
var routine_wg *sync.WaitGroup
var max_routine int32

func init() {
	func_chan = make(chan func())
	atomic.StoreInt32(&routine_count, 0)
	routine_wg = &sync.WaitGroup{}
}

func Init(maxroutine int) {
	max_routine = int32(maxroutine)

	go func() {
		for {
			select {
			case f := <-func_chan:
				routine_wg.Wait()

				if atomic.AddInt32(&routine_count, 1) == max_routine {
					routine_wg.Add(1)
				}

				go func() {
					f()
					if atomic.AddInt32(&routine_count, -1) == max_routine-1 {
						routine_wg.Done()
					}
				}()
			}
		}
	}()
}

func BlanceRun(f func()) {
	func_chan <- f
}
