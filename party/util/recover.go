package util

import (
	"fmt"

	"github.com/go-errors/errors"
)

// Protect ...
func Protect(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("-------------------Protect error-------------------")
			fmt.Println("Recovery_error: ", errors.Wrap(err, 2).ErrorStack())
		}
	}()
	f()
}
