package form

type LoginReq struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id"`
}

type LoginResp struct {
	AssessToken           string `json:"assess_token"`
	AssessTokenExpiresIn  int64  `json:"assess_token_expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	Uid                   uint64 `json:"uid"`
}
