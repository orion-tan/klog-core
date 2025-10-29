package utils

import "regexp"

// 验证是否为正确邮箱格式
// @email 邮箱
// @return 是否为正确邮箱格式
func ValidateEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}

// 验证是否为正确 URL 格式
// @url URL
// @return 是否为正确 URL 格式
func ValidateURL(url string) bool {
	return regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(url)
}

// 验证是否为正确用户名
// @username 用户名
// @return 是否为正确用户名
func ValidateUsername(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+$`).MatchString(username)
}

// 验证是否为正确密码格式
// @password 密码
// @return 是否为正确密码格式
func ValidatePassword(password string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+$`).MatchString(password)
}
