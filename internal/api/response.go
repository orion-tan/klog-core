package api

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string `json:"token"`
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}
