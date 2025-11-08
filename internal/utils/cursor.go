package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CursorData 游标数据结构
type CursorData struct {
	SortField string // 排序字段名
	SortValue string // 排序字段值（字符串形式）
	ID        uint   // 记录ID（用于唯一性）
}

// EncodeCursor 编码游标
// 格式：base64(sortField:sortValue:id)
func EncodeCursor(data CursorData) string {
	raw := fmt.Sprintf("%s:%s:%d", data.SortField, data.SortValue, data.ID)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor 解码游标
func DecodeCursor(cursor string) (*CursorData, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("无效的游标格式: %w", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("游标格式错误")
	}

	id, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的ID: %w", err)
	}

	return &CursorData{
		SortField: parts[0],
		SortValue: parts[1],
		ID:        uint(id),
	}, nil
}

// FormatSortValue 格式化排序字段值为字符串（用于游标）
func FormatSortValue(field string, value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case *time.Time:
		if v == nil {
			return ""
		}
		return v.Format(time.RFC3339Nano)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%v", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ParseSortValue 解析游标中的排序值为对应类型
func ParseSortValue(field string, valueStr string) (interface{}, error) {
	if valueStr == "" {
		return nil, nil
	}

	// 根据常见字段名推断类型
	switch field {
	case "published_at", "created_at", "updated_at":
		return time.Parse(time.RFC3339Nano, valueStr)
	case "view_count":
		return strconv.ParseUint(valueStr, 10, 64)
	case "title", "slug":
		return valueStr, nil
	default:
		// 默认尝试解析为数字，失败则作为字符串
		if num, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return num, nil
		}
		return valueStr, nil
	}
}
