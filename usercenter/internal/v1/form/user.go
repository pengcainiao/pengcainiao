package form

type RegisterUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterAndLoginResponse struct {
	UserID string `json:"user_id,omitempty" db:"id"` //用户ID
}

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
