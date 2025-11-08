package api

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// UserGetMeResponse 用户获取自身信息响应
type UserGetMeResponse struct {
	ID        uint    `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Nickname  string  `json:"nickname"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string `json:"token"`
}

// PaginatedResponse 分页响应（传统offset分页）
type PaginatedResponse struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

// CursorPaginatedResponse 游标分页响应
type CursorPaginatedResponse struct {
	Data       interface{} `json:"data"`
	NextCursor *string     `json:"next_cursor"` // 下一页游标，为nil表示没有下一页
	PrevCursor *string     `json:"prev_cursor"` // 上一页游标，为nil表示没有上一页
	HasMore    bool        `json:"has_more"`    // 是否有更多数据
	Limit      int         `json:"limit"`       // 当前每页数量
}

// CategoryWithCount 分类及其文章数量
type CategoryWithCount struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	PostCount   int64   `json:"post_count"`
}

// TagWithCount 标签及其文章数量
type TagWithCount struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PostCount int64  `json:"post_count"`
}

// CommentResponse 评论响应（IP已脱敏）
type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    *uint     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	IP        string    `json:"ip"` // 脱敏后的IP
	Status    string    `json:"status"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt string    `json:"created_at"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}
