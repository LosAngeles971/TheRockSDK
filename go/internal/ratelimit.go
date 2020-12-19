package internal

import (
	"time"
)

var firstcall int64 = 0
var lastcall int64 = 0
var count int = 0

// wait time if you reached a rate limit
func WaitDueRateLimit(reqPerSecond int) {
	count++ 
	// timestamp with millisecond precision
	lastcall = time.Now().UnixNano() / 1000000
	if count == 1 || firstcall == 0 {
		// this is the really first call
		firstcall = lastcall
		return
	}
	d := lastcall - firstcall
	if d > 1000 {
		// if delay between last and first is greater than 1 second, start a new cycle of counting
		count = 1
		firstcall = lastcall
		return
	}
	if count <= reqPerSecond {
		//limit not reached yet, so no sleep and calls at the maximum speed
		return
	}
	sleep_time := time.Duration(firstcall + 1000 - lastcall)
	time.Sleep(sleep_time * time.Millisecond)
	// start a new cycle of counting
	count = 0
}