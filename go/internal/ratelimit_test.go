package internal

import (
	"time"
	"fmt"
	"testing"
)

func TestRateLimit(t *testing.T) {
	start := time.Now()
	before := time.Now()
    for i := 0; i < 10; i++ {
		WaitDueRateLimit(5)
		now := time.Now()
        fmt.Println(i, now.Sub(before), now.Sub(start))
        before = now
	}
}