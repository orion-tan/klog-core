package model

import (
	"time"
)

// User 用户表（单用户设计，仅存储唯一管理员账号）
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // 不在 JSON 中返回密码
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Nickname  string    `gorm:"type:varchar(50);not null" json:"nickname"`
	Bio       *string   `gorm:"type:text" json:"bio,omitempty"`
	AvatarURL *string   `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 关联关系
	Posts    []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// Post 文章表
type Post struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID    *uint      `gorm:"index" json:"category_id"` // 可以没有分类
	AuthorID      uint       `gorm:"not null;index" json:"author_id"`
	Title         string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug          string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Content       string     `gorm:"type:longtext;not null" json:"content"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`
	CoverImageURL string     `gorm:"type:varchar(255)" json:"cover_image_url"`
	Status        string     `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"` // 'draft', 'published', 'archived'
	ViewCount     uint64     `gorm:"not null;default:0" json:"view_count"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 关联关系
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:RESTRICT" json:"author,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// Category 分类表
type Category struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug        string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Description *string `gorm:"type:text" json:"description"`

	// 关联关系
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// Tag 标签表
type Tag struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug string `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`

	// 关联关系
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// PostTag 文章<->标签关联表（多对多中间表）
type PostTag struct {
	PostID uint `gorm:"primaryKey;constraint:OnDelete:CASCADE" json:"post_id"`
	TagID  uint `gorm:"primaryKey;constraint:OnDelete:CASCADE" json:"tag_id"`
}

// Comment 评论表
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	UserID    *uint     `gorm:"index" json:"user_id"`           // 可为空（游客评论）
	Name      string    `gorm:"type:varchar(50)" json:"name"`   // 游客评论时使用
	Email     string    `gorm:"type:varchar(100)" json:"email"` // 游客评论时使用
	Content   string    `gorm:"type:text;not null" json:"content"`
	IP        string    `gorm:"type:varchar(100);not null" json:"ip"`
	Status    string    `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"` // 'pending', 'approved', 'spam'
	ParentID  *uint     `gorm:"index" json:"parent_id"`                                          // 可为空（顶级评论）
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`

	// 关联关系
	Post    Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	User    *User     `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"replies,omitempty"`
}

// Setting 系统设置表
type Setting struct {
	Key   string `gorm:"type:varchar(50);primaryKey" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Type  string `gorm:"type:varchar(20);not null;default:'str'" json:"type"` // 'str', 'number', 'json'
}

// Media 媒体库表（单用户设计，无需记录上传者信息）
type Media struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath  string    `gorm:"type:varchar(255);not null" json:"file_path"`
	FileHash  string    `gorm:"type:varchar(255);index" json:"file_hash"`
	MimeType  string    `gorm:"type:varchar(100);not null" json:"mime_type"`
	Size      int64     `gorm:"not null" json:"size"` // 文件大小（字节）
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

// 指定表名
func (User) TableName() string {
	return "users"
}

func (Post) TableName() string {
	return "posts"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (PostTag) TableName() string {
	return "post_tags"
}

func (Comment) TableName() string {
	return "comments"
}

func (Setting) TableName() string {
	return "settings"
}

func (Media) TableName() string {
	return "media"
}
