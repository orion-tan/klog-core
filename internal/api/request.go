package api

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=30"`
	Nickname string `json:"nickname" binding:"required,min=3,max=50"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=30"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Nickname    *string `json:"nickname,omitempty"`
	Username    *string `json:"username,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Email       *string `json:"email,omitempty"`
	OldPassword *string `json:"old_password,omitempty"`
	NewPassword *string `json:"new_password,omitempty"`
}

// PostCreateRequest 创建文章请求
type PostCreateRequest struct {
	CategoryID    *uint    `json:"category_id"`
	Title         string   `json:"title" binding:"required"`
	Slug          string   `json:"slug" binding:"required"`
	Content       string   `json:"content" binding:"required"`
	Excerpt       string   `json:"excerpt"`
	CoverImageURL string   `json:"cover_image_url"`
	Status        string   `json:"status" binding:"required,oneof=draft published archived"`
	Tags          []string `json:"tags"`
}

// PostUpdateRequest 更新文章请求
type PostUpdateRequest struct {
	CategoryID    *uint     `json:"category_id,omitempty"`
	Title         *string   `json:"title,omitempty"`
	Slug          *string   `json:"slug,omitempty"`
	Content       *string   `json:"content,omitempty"`
	Excerpt       *string   `json:"excerpt,omitempty"`
	CoverImageURL *string   `json:"cover_image_url,omitempty"`
	Status        *string   `json:"status" binding:"omitempty,oneof=draft published archived"`
	Tags          *[]string `json:"tags,omitempty"`
}

// CategoryCreateRequest 创建分类请求
type CategoryCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

// CategoryUpdateRequest 更新分类请求
type CategoryUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// CommentCreateRequest 创建评论请求
type CommentCreateRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// CommentUpdateRequest 更新评论请求
type CommentUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved spam"`
}

// SettingUpsertRequest 创建/更新设置请求
type SettingUpsertRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	Type  string `json:"type" binding:"required,oneof=str number json"`
}

// SettingBatchUpsertRequest 批量创建/更新设置请求
type SettingBatchUpsertRequest struct {
	Settings []SettingUpsertRequest `json:"settings" binding:"required,min=1"`
}
