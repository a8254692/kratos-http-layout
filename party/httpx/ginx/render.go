package ginx

import (
	"io/ioutil"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	jsoniter "github.com/json-iterator/go"

	"gitlab.top.slotssprite.com/my/api-layout/party/httpx"
	"gitlab.top.slotssprite.com/my/api-layout/party/runtimex"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

// Render gin render
func (c *Context) Render(status statusx.Status, data interface{}, httpCode ...int) {
	result := httpx.NewResult(status, data)
	result.RequestId = tracing.TraceID()(c.Request.Context()).(string)

	c.Header("Cache-Control", "no-cache")
	_code := http.StatusOK
	if len(httpCode) > 0 {
		_code = httpCode[0]
		c.Status(_code)
	}
	c.log(result)
	c.JSON(_code, result)
}

// RenderWithMsg ...
func (c *Context) RenderWithMsg(status statusx.Status, data interface{}, msg string, httpCode ...int) {
	result := httpx.NewResultWithMsg(status, data, msg)
	result.RequestId = tracing.TraceID()(c.Request.Context()).(string)

	c.Header("Cache-Control", "no-cache")
	_code := http.StatusOK
	if len(httpCode) > 0 {
		_code = httpCode[0]
		c.Status(_code)
	}
	c.log(result)
	c.JSON(_code, result)
}

// RenderResult ...
func (c *Context) RenderResult(result *httpx.Result, httpCode ...int) {
	if result != nil {
		result.RequestId = tracing.TraceID()(c.Request.Context()).(string)
	}

	c.Header("Cache-Control", "no-cache")
	_code := http.StatusOK
	if len(httpCode) > 0 {
		_code = httpCode[0]
		c.Status(_code)
	}
	c.log(result)
	c.JSON(_code, result)
}

// RenderText ...
func (c *Context) RenderText(code int, text string) {
	result := &httpx.Result{Data: text}

	c.Header("Cache-Control", "no-cache")
	c.Status(code)
	c.log(result)
	c.String(code, text)
}

// log ...
func (c *Context) log(result *httpx.Result) {
	if result == nil {
		return
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	urlParams, _ := jsoniter.MarshalToString(c.Request.URL.Query())
	//headers, _ := jsoniter.MarshalToString(map[string]string{
	//	"X-Channel": c.GetHeader("X-Channel"),
	//	"X-Agent":   c.GetHeader("X-AGENT"),
	//})
	params := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.FullPath(),
		"param":      c.paramsString(),
		"url_params": urlParams,
		//"headers":   headers,
		"body": string(body),
	}

	if metadata := c.GetStringMap("metadata"); metadata != nil {
		if userId, ok := metadata["user_id"]; ok {
			params["user_id"] = userId
		}
	}

	_result := map[string]interface{}{
		"status":    result.Status,
		"message":   result.Msg,
		"code":      result.Code,
		"http_code": c.Writer.Status(),
	}
	if _, ok := result.Data.(string); ok {
		_result["data"] = result.Data
	}
	if err, ok := result.Data.(error); ok {
		result.Data = nil
		_result["error"] = err.Error()
	}

	log.Context(c.Request.Context()).Log(log.LevelWarn,
		"@render", runtimex.Caller(3)(c.Request.Context()),
		"@errorf", result.Errorf,
		"@field", map[string]interface{}{
			"params":       params,
			"process_time": c.processTime(),
			"result":       _result,
		},
	)
}
