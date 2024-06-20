package status

import (
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	pbresult "gitlab.top.slotssprite.com/my/api-layout/api/helloworld/v1/result"
	"gitlab.top.slotssprite.com/my/api-layout/party/httpx"
	"gitlab.top.slotssprite.com/my/api-layout/party/statusx"
)

const (
	// BusinessError 自定义的业务错误 code
	BusinessError codes.Code = 101

	// _code HACK
	_code = "11000001"
)

// NewError ...
func NewError(err error, status statusx.Status) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	st, _ := gstatus.New(BusinessError, errMsg).WithDetails(&pbresult.Result{
		Status: status.String(),
		Msg:    statusx.GetMsg(status),
		Code:   _code,
	})
	return st.Err()
}

// FromError ...
func FromError(err error) *httpx.Result {
	return fromError(err, false)
}

// FromErrorWithMsg ...
func FromErrorWithMsg(err error) *httpx.Result {
	return fromError(err, true)
}

// IsNotFoundError ...
func IsNotFoundError(result *httpx.Result) bool {
	if result == nil {
		return false
	}
	return result.Status == statusx.StatusNotFound
}

// fromError ...
func fromError(err error, useErrorMsg bool) *httpx.Result {
	if err == nil {
		return nil
	}
	st, ok := gstatus.FromError(err)
	if !ok {
		return nil
	}
	switch st.Code() {
	case codes.Unknown:
		return httpx.FromErrorResult(statusx.StatusUnknownError, err)
	case codes.InvalidArgument:
		return httpx.FromErrorResult(statusx.StatusInvalidRequest, err)
	case codes.DeadlineExceeded:
		return httpx.FromErrorResult(statusx.StatusRequestTimeout, err)
	case BusinessError:
		var (
			_result *pbresult.Result
			_ok     bool
		)
		for _, detail := range st.Details() {
			if _result, _ok = detail.(*pbresult.Result); _ok {
				break
			}
		}
		if _result != nil {
			if useErrorMsg {
				return httpx.ConvertResultWithMsg(_result, err)
			}
			return httpx.ConvertResult(_result, err)
		}
		return httpx.FromErrorResult(statusx.StatusInternalServerError, err)
	default:
		return httpx.FromErrorResult(statusx.StatusInternalServerError, err)
	}
}
