package form

type RegisterAndLoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterAndLoginResponse struct {
	UserID string `json:"user_id,omitempty" db:"id"` //用户ID
	Token  string `json:"Token,omitempty"`           //用户token，仅用于在登录时与用户信息一起返回token
}
