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

// ValidateMarkdownContent 验证Markdown内容是否安全
// @content Markdown内容
// @return 是否为安全的Markdown内容
func ValidateMarkdownContent(content string) bool {
	if content == "" {
		return false
	}

	// 检测危险的HTML标签和JavaScript
	dangerousPatterns := []string{
		`<script[\s\S]*?>[\s\S]*?</script>`,                    // <script>标签
		`<iframe[\s\S]*?>`,                                      // <iframe>标签
		`javascript:`,                                           // javascript:协议
		`on\w+\s*=`,                                             // onclick等事件处理器
		`<embed[\s\S]*?>`,                                       // <embed>标签
		`<object[\s\S]*?>`,                                      // <object>标签
		`<applet[\s\S]*?>`,                                      // <applet>标签
		`<meta[\s\S]*?>`,                                        // <meta>标签
		`<link[\s\S]*?>`,                                        // <link>标签
		`<base[\s\S]*?>`,                                        // <base>标签
		`<form[\s\S]*?>`,                                        // <form>标签
		`data:text/html`,                                        // data URI
		`vbscript:`,                                             // vbscript:协议
	}

	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(`(?i)`+pattern, content)
		if matched {
			return false
		}
	}

	return true
}

// SanitizeMarkdownContent 清理Markdown内容（移除危险字符但保留内容）
// @content Markdown内容
// @return 清理后的内容
func SanitizeMarkdownContent(content string) string {
	// 移除潜在的XSS攻击向量
	sanitized := content

	// 移除<script>标签及其内容
	sanitized = regexp.MustCompile(`(?i)<script[\s\S]*?>[\s\S]*?</script>`).ReplaceAllString(sanitized, "")

	// 移除危险的HTML标签
	dangerousTags := []string{"iframe", "embed", "object", "applet", "meta", "link", "base", "form"}
	for _, tag := range dangerousTags {
		pattern := `(?i)<` + tag + `[\s\S]*?>`
		sanitized = regexp.MustCompile(pattern).ReplaceAllString(sanitized, "")
	}

	// 移除事件处理器属性
	sanitized = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`).ReplaceAllString(sanitized, "")

	// 移除javascript:和data:协议
	sanitized = regexp.MustCompile(`(?i)javascript:`).ReplaceAllString(sanitized, "")
	sanitized = regexp.MustCompile(`(?i)data:text/html`).ReplaceAllString(sanitized, "")

	return sanitized
}

// MaskIP 对IP地址进行脱敏处理
// @ip 原始IP地址
// @return 脱敏后的IP地址
func MaskIP(ip string) string {
	if ip == "" {
		return ""
	}

	// IPv4 脱敏：保留前两段，后两段替换为 *
	// 例如：192.168.1.1 -> 192.168.*.*
	ipv4Pattern := regexp.MustCompile(`^(\d{1,3}\.\d{1,3})\.\d{1,3}\.\d{1,3}$`)
	if ipv4Pattern.MatchString(ip) {
		matches := ipv4Pattern.FindStringSubmatch(ip)
		if len(matches) > 1 {
			return matches[1] + ".*.*"
		}
	}

	// IPv6 脱敏：保留前两组，其余替换为 *
	// 例如：2001:0db8:85a3:0000:0000:8a2e:0370:7334 -> 2001:0db8:*:*:*:*:*:*
	ipv6Pattern := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:[0-9a-fA-F]{1,4}):.*$`)
	if ipv6Pattern.MatchString(ip) {
		matches := ipv6Pattern.FindStringSubmatch(ip)
		if len(matches) > 1 {
			return matches[1] + ":*:*:*:*:*:*"
		}
	}

	// IPv6 简写形式（包含::）的脱敏
	ipv6ShortPattern := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:[0-9a-fA-F]{1,4}):.*$`)
	if ipv6ShortPattern.MatchString(ip) {
		matches := ipv6ShortPattern.FindStringSubmatch(ip)
		if len(matches) > 1 {
			return matches[1] + "::*"
		}
	}

	// 如果无法识别格式，返回通用脱敏
	return "***"
}
