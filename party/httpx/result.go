package httpx

import (
	"context"

	pbresult "gitlab.top.slotssprite.com/my/api-layout/api/helloworld/v1/result"
	"gitlab.top.slotssprite.com/my/api-layout/party/runtimex"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

// Result ...
type Result struct {
	Status    statusx.Status `json:"status"`
	Msg       string         `json:"message"`
	Data      interface{}    `json:"data"`
	Code      string         `json:"code"`
	RequestId string         `json:"requestId"`
	Errorf    string         `json:"-"`
}

// NewResult use default msg
func NewResult(status statusx.Status, data interface{}) *Result {
	return &Result{Status: status, Msg: statusx.GetMsg(status), Data: data, Code: _code(status), Errorf: runtimex.Caller(2)(context.Background()).(string)}
}

// NewResultWithMsg use custom msg
func NewResultWithMsg(status statusx.Status, data interface{}, msg string) *Result {
	return &Result{Status: status, Msg: msg, Data: data, Code: _code(status), Errorf: runtimex.Caller(2)(context.Background()).(string)}
}

// FromErrorResult use for status.fromError
func FromErrorResult(status statusx.Status, data interface{}) *Result {
	return &Result{Status: status, Msg: statusx.GetMsg(status), Data: data, Code: _code(status), Errorf: runtimex.Caller(4)(context.Background()).(string)}
}

// ConvertResult use for status.fromError
func ConvertResult(result *pbresult.Result, err error) *Result {
	status := statusx.Status(result.Status)
	return &Result{Status: status, Msg: statusx.GetMsg(status), Data: err, Code: _code(status), Errorf: runtimex.Caller(4)(context.Background()).(string)}
}

// ConvertResultWithMsg use for status.fromError
func ConvertResultWithMsg(result *pbresult.Result, err error) *Result {
	status := statusx.Status(result.Status)
	return &Result{Status: status, Msg: result.Msg, Data: err, Code: _code(status), Errorf: runtimex.Caller(4)(context.Background()).(string)}
}

// _code OPTIMIZE: code 做兼容
func _code(status statusx.Status) string {
	switch status {
	case statusx.StatusOk:
		return "12000000"
	default:
		return "12000001"
	}
}
