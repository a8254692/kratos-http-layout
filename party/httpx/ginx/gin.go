package ginx

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Context ...
type Context struct {
	*gin.Context
}

// NewContext ...
func NewContext(ctx *gin.Context) *Context {
	return &Context{ctx}
}

// Handle ...
func Handle(f func(ctx *Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(NewContext(c))
	}
}

// paramsString ...
func (c *Context) paramsString() string {
	var paramsBuff bytes.Buffer
	for _, param := range c.Params {
		paramsBuff.WriteString(param.Key)
		paramsBuff.WriteString("=")
		paramsBuff.WriteString(param.Value)
		paramsBuff.WriteString(",")
	}
	if paramsBuff.String() != "" {
		return strings.TrimSuffix(paramsBuff.String(), ",")
	}
	return paramsBuff.String()
}

// processTime ...
func (c *Context) processTime() (processTime int64) {
	startAt := c.GetInt64("startAt")
	if startAt > 0 {
		processTime = time.Now().UnixMilli() - startAt
	}
	return
}

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}
