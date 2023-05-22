package models

import "errors"

var 	UserNotExistsError = errors.New("user not exists")

const (
	//用户在线设备列表，每次访问均会更新过期时间
	// u:onlie:{user_id}:{device_id}
	// HASH
	//{device_id}:{Token}
	RedisKeyForOnlineUsers = `u:online:%s`
)

