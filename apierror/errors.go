package apierror

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error API 返回错误的基本结构，以前用 ApiError 这个玩意限制的太死了，小弟改成了接口。
type Error interface {
	Code() int
	Message() string
}

var _ Error = (*ApiError)(nil)

// ApiError 建议直接使用 NewApiError 或者 NewApiErrWithData 方法，避免后期难以维护。
type ApiError struct {
	code    int
	message string
}

var (
	InvalidParamErr = &ApiError{
		code:    1,
		message: "无效的参数",
	}
	DeleteBannerErr = &ApiError{
		code:    1,
		message: "一个场景下至少要有一个banner",
	}
	InternalErr = &ApiError{
		code:    2,
		message: "服务器内部错误",
	}
	LocRechargeErr = &ApiError{
		code:    2,
		message: "充值操作太快了",
	}
	DuplicateGroupErr = &ApiError{
		code:    3,
		message: "不能多次设置同一个群",
	}
	TooMuchDataErr = &ApiError{
		code:    5,
		message: "禁止请求",
	}
	IndexExistsErr = &ApiError{
		code:    4,
		message: "列表位置已存在",
	}

	// 1** 无效的参数
	InvalidTimeParamErr = &ApiError{
		code:    101,
		message: "相同游戏相同地区开始和结束时间不能交叉",
	}
	InvalidUidParamErr = &ApiError{
		code:    102,
		message: "存在无效的uid",
	}
	InvalidPastTimeParamErr = &ApiError{
		code:    103,
		message: "不得使用过去时间",
	}
	InvalidUserErr = &ApiError{
		code:    104,
		message: "无效的用户",
	}
	InvalidTimeErr = &ApiError{
		code:    105,
		message: "无效的时间",
	}
	InvalidSlice = &ApiError{
		code:    106,
		message: "无效的数组",
	}
	InvalidUser = &ApiError{
		code:    107,
		message: "賬號被凍結",
	}
	CoinLimitSecErr = &ApiError{
		code:    109,
		message: "今日代儲已達上限，請聯繫一級代理",
	}
	CoinLimitFirErr = &ApiError{
		code:    119,
		message: "賬戶金幣餘額不足，請聯繫官方",
	}
	LoginTimesOneHourLimitFErr = &ApiError{
		code:    120,
		message: "登錄失敗次數過多，請1小時後再試",
	}
	BossUidRechargeLimit = &ApiError{
		code:    121,
		message: "代儲失敗，對方被限制充值",
	}

	UnauthorizedErr = &ApiError{
		code:    401,
		message: "无效的身份认证",
	}
	ForbiddenErr = &ApiError{
		code:    403,
		message: "没有权限",
	}
	NilTargetResultErr = &ApiError{
		code:    404,
		message: "目标参数不存在",
	}
	BannerLenLimitErr = &ApiError{
		code:    404,
		message: "该场景下banner数量已到达上限",
	}
	NameExistErr = &ApiError{
		code:    410,
		message: "名字已存在",
	}
	AlreadyDeletedErr = &ApiError{
		code:    411,
		message: "已删除",
	}
	NotFoundErr = &ApiError{
		code:    404,
		message: "未找到",
	}
	ItemNotExistErr = &ApiError{
		code:    404,
		message: "业务ID填写错误",
	}
	UserNotExistErr = &ApiError{
		code:    404,
		message: "用户不存在",
	}


	//
)

func (e *ApiError) Error() string {
	return fmt.Sprintf("err code %d with message %s", e.code, e.message)
}

func (e *ApiError) Code() int {
	return e.code
}

func (e *ApiError) Message() string {
	return e.message
}

func Fail(c *gin.Context, err Error) {
	data := gin.H{"code": err.Code(), "msg": err.Message()}
	if v, ok := err.(ErrorWithData); ok {
		data["data"] = v.Data()
	}
	c.AbortWithStatusJSON(http.StatusOK, data)
}

func FailWithMsg(ctx *gin.Context, err *ApiError, msg string) {
	Fail(ctx, NewApiError(err.Code(), msg))
}

func NewApiError(code int, msg string) Error {
	return &ApiError{
		code:    code,
		message: msg,
	}
}

func NewGeneralError(errArg error) Error {
	return InternalErr
	//marshal, err := json.Marshal(&errArg)
	//if err != nil {
	//	return InternalErr
	//}
	//
	//log.Debugf("error value: %+v; error text: %s", errArg, string(marshal))
	//
	//grpcErr := &struct {
	//	code int
	//	msg  string
	//}{}
	//err = json.Unmarshal(marshal, grpcErr)
	//
	//if err != nil {
	//	return InternalErr
	//}
	//
	//return &ApiError{
	//	code:    grpcErr.code,
	//	message: grpcErr.msg,
	//}
}
