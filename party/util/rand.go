package util

import (
	"math/rand"
	"time"
)

var (
	// randSrc 伪随机source, 供 util 包使用
	randSrc = rand.NewSource(time.Now().UnixNano())
)
