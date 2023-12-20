package sms

import (
	"sync/atomic"
	"time"
)

// TODO: Remove this file, this file only for temporary development!

var fakeMessageID atomic.Int64
var fakeBatchID atomic.Int64

func init() {
	fakeMessageID.Store(time.Now().UnixNano())
	fakeBatchID.Store(time.Now().UnixNano())
}
