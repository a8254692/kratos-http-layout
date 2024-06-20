package runtimex

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-kratos/kratos/v2/log"
)

// Caller ...
func Caller(depth int) log.Valuer {
	return func(context.Context) interface{} {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			return ""
		}

		var count int
		for i := len(file) - 1; i >= 0; i-- {
			if file[i] == '/' {
				count++
			}

			if count == 3 {
				return fmt.Sprintf("%s:%d", file[i+1:], line)
			}
		}
		return fmt.Sprintf("%s:%d", file, line)
	}
}
