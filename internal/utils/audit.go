package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLog 审计日志结构
type AuditLog struct {
	Timestamp  time.Time `json:"timestamp"`
	UserID     uint      `json:"user_id,omitempty"`
	Username   string    `json:"username,omitempty"`
	IP         string    `json:"ip"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource,omitempty"`
	ResourceID string    `json:"resource_id,omitempty"`
	Details    string    `json:"details,omitempty"`
	Result     string    `json:"result"` // success, failure
}

// LogAudit 记录审计日志
// @c Gin上下文
// @action 操作类型（如 "login", "upload_file", "delete_file"）
// @resource 资源类型（如 "media", "post", "user"）
// @resourceID 资源ID
// @details 详细信息
// @result 操作结果（"success" 或 "failure"）
func LogAudit(c *gin.Context, action, resource, resourceID, details, result string) {
	auditLog := AuditLog{
		Timestamp:  time.Now(),
		IP:         c.ClientIP(),
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		Result:     result,
	}

	// 获取用户信息（如果已认证）
	if claims, exists := c.Get("claims"); exists {
		if klogClaims, ok := claims.(*KLogClaims); ok {
			auditLog.UserID = klogClaims.UserID
			auditLog.Username = klogClaims.Username
		}
	}

	// 使用结构化日志记录
	SugarLogger.Infow("审计日志",
		"timestamp", auditLog.Timestamp,
		"user_id", auditLog.UserID,
		"username", auditLog.Username,
		"ip", auditLog.IP,
		"action", auditLog.Action,
		"resource", auditLog.Resource,
		"resource_id", auditLog.ResourceID,
		"details", auditLog.Details,
		"result", auditLog.Result,
	)
}

// 记录注册日志
func LogRegister(c *gin.Context, username string, success bool) {
	result := "success"
	details := "用户注册成功"
	if !success {
		result = "failure"
		details = "用户注册失败"
	}
	LogAudit(c, "register", "auth", username, details, result)
}

// LogLogin 记录登录日志
func LogLogin(c *gin.Context, username string, success bool) {
	result := "success"
	details := "用户登录成功"
	if !success {
		result = "failure"
		details = "用户登录失败"
	}
	LogAudit(c, "login", "auth", username, details, result)
}

// LogFileOperation 记录文件操作日志
func LogFileOperation(c *gin.Context, operation, fileID, fileName string, success bool) {
	result := "success"
	details := fileName
	if !success {
		result = "failure"
	}
	LogAudit(c, operation, "media", fileID, details, result)
}

// LogPermissionChange 记录权限变更日志
func LogPermissionChange(c *gin.Context, targetUserID, details string, success bool) {
	result := "success"
	if !success {
		result = "failure"
	}
	LogAudit(c, "permission_change", "user", targetUserID, details, result)
}
