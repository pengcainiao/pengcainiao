package httprouter

import "errors"

const (
	ErrOperateFailed            = 40001 // 发生错误，但可忽略
	ErrInvalidParameterCode     = 40002 // 参数无效
	ErrInternalErrorCode        = 40006 // 服务器错误
	ErrUnAuthorizedCode         = 40008 // 无权限访问
	ErrAliyunServiceCode        = 40009 // 阿里云服务异常
	ErrNotFoundCode             = 40010 // 访问数据不存在
	ErrQuantityOutOfLimitCode   = 40011 // 数量超出限制
	ErrInvalidTaskStateCode     = 40012 // 无效任务状态
	ErrNoEnoughSpaceCode        = 40013 // 可用空间不足
	ErrAlreadyExistsCode        = 40017 // 记录已存在
	ErrVoteCancelledCode        = 40018 // 投票已取消
	ErrVoteFinishedCode         = 40019 // 投票已结束
	ErrVoteCompletedCode        = 40020 // 投票已完成
	ErrTaskCancelledCode        = 40021 // 任务已取消
	ErrFreqOutOfLimitCode       = 40022 // 请求频率超过限制
	ErrTaskWithdrawnCode        = 40023 // 任务已撤回
	ErrRefTaskOutOfLimitCode    = 40024 // 笔记关联事项超出限制
	ErrRiskyContentCode         = 40025 // 包含敏感内容
	ErrTaskUnCompleteCode       = 40027 // 任务未完成
	ErrTaskDeletedCode          = 40028 // 任务已删除
	ErrAccountAlreadyBindedCode = 40031 // 账号已经被绑定过了
	ErrLoginFailedCode          = 40032 // 登录失败
	ErrOtherPlatformExecuteCode = 40036 // 其他平台操作过了
	ErrCacheExpiredCode         = 40037 // 缓存已过期
	ErrEditingCode              = 40038 // 正在编辑
	ErrLogoutCode               = 40044 // 注销了小程序登录
	ErrReqReleaseAuthCode       = 40045 // 没有权限访问测试环境
	ErrUpgradeVersionCode       = 40048 // 需要升级版本
	ErrInvalidTokenCode         = 40049 // 无效的token，需要退登
	ErrCustomMessageCode        = 40050 // 自定义 message code，适用于前端直接展示 message 信息
	ErrUserRemovedCode          = 40051 // 用户已被移出
	ErrVersionExpiredCode       = 40052 // 版本已到期
	ErrUnBindYzjAccountCode     = 40053 // 未绑定云之家账号
	ErrUpgradeInProgressCode    = 49999 // 服务器正在升级
)

var (
	err40001 = &Response{Code: ErrOperateFailed, Message: "操作失败"}
	err40002 = &Response{Code: ErrInvalidParameterCode, Message: "输入参数不完整或无效"}
	err40006 = &Response{Code: ErrInternalErrorCode, Message: "请求错误"}
	err40008 = &Response{Code: ErrUnAuthorizedCode, Message: "无权限访问"}
	err40009 = &Response{Code: ErrAliyunServiceCode, Message: "阿里云服务异常"}
	err40010 = &Response{Code: ErrNotFoundCode, Message: "访问数据不存在"}
	err40011 = &Response{Code: ErrQuantityOutOfLimitCode, Message: "数量超出限制"}
	err40012 = &Response{Code: ErrInvalidTaskStateCode, Message: "无效任务状态"}
	err40013 = &Response{Code: ErrNoEnoughSpaceCode, Message: "可用空间不足"}
	err40017 = &Response{Code: ErrAlreadyExistsCode, Message: "记录已存在"}
	err40018 = &Response{Code: ErrVoteCancelledCode, Message: "投票已取消"}
	err40019 = &Response{Code: ErrVoteFinishedCode, Message: "投票已结束"}
	err40020 = &Response{Code: ErrVoteCompletedCode, Message: "投票已完成"}
	err40021 = &Response{Code: ErrTaskCancelledCode, Message: "任务已取消"}
	err40022 = &Response{Code: ErrFreqOutOfLimitCode, Message: "请求频率超过限制"}
	err40023 = &Response{Code: ErrTaskWithdrawnCode, Message: "任务已撤回"}
	err40024 = &Response{Code: ErrRefTaskOutOfLimitCode, Message: "笔记关联事项超出限制"}
	err40025 = &Response{Code: ErrRiskyContentCode, Message: "包含敏感内容"}
	err40027 = &Response{Code: ErrTaskUnCompleteCode, Message: "任务未完成"}
	err40028 = &Response{Code: ErrTaskDeletedCode, Message: "任务已删除"}
	err40031 = &Response{Code: ErrAccountAlreadyBindedCode, Message: "账号已经被绑定过了"}
	err40032 = &Response{Code: ErrLoginFailedCode, Message: "登录失败"}
	err40036 = &Response{Code: ErrOtherPlatformExecuteCode, Message: "其他平台操作过了"}
	err40037 = &Response{Code: ErrCacheExpiredCode, Message: "缓存已过期"}
	err40038 = &Response{Code: ErrEditingCode, Message: "正在编辑"}
	err40048 = &Response{Code: ErrUpgradeVersionCode, Message: "需要升级版本"}
	err40049 = &Response{Code: ErrInvalidTokenCode, Message: "无效的token"}
	err40050 = &Response{Code: ErrCustomMessageCode, Message: "自定义 message"}
	err40051 = &Response{Code: ErrUserRemovedCode, Message: "用户已被移出"}
	err40052 = &Response{Code: ErrVersionExpiredCode, Message: "版本已到期"}
	err40053 = &Response{Code: ErrUnBindYzjAccountCode, Message: "未绑定云之家账号"}
	err49999 = &Response{Code: ErrUpgradeInProgressCode, Message: "服务器正在升级"}

	errorConfig = map[int]*Response{
		err40001.Code: err40001,
		err40002.Code: err40002,
		err40006.Code: err40006,
		err40008.Code: err40008,
		err40009.Code: err40009,
		err40010.Code: err40010,
		err40011.Code: err40011,
		err40012.Code: err40012,
		err40013.Code: err40013,
		err40017.Code: err40017,
		err40018.Code: err40018,
		err40019.Code: err40019,
		err40020.Code: err40020,
		err40021.Code: err40021,
		err40022.Code: err40022,
		err40023.Code: err40023,
		err40024.Code: err40024,
		err40025.Code: err40025,
		err40027.Code: err40027,
		err40028.Code: err40028,
		err40031.Code: err40031,
		err40032.Code: err40032,
		err40036.Code: err40036,
		err40037.Code: err40037,
		err40038.Code: err40038,
		err40048.Code: err40048,
		err40049.Code: err40049,
		err40050.Code: err40050,
		err40051.Code: err40051,
		err40052.Code: err40052,
		err40053.Code: err40053,
		err49999.Code: err49999,
	}
)

func GetError(code int, internalError error) Response {
	if v, ok := errorConfig[code]; ok {
		if internalError != nil {
			v.InternalError = internalError.Error()
		}
		return *v
	}
	if internalError == nil {
		internalError = errors.New("internal error is Nil")
	}
	return Response{
		Code:          code,
		InternalError: internalError.Error(),
	}
}
