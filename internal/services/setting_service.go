package services

import (
	"encoding/json"
	"errors"
	"klog-backend/internal/api"
	"klog-backend/internal/model"
	"klog-backend/internal/repository"
	"strconv"
)

type SettingService struct {
	settingRepo *repository.SettingRepository
}

func NewSettingService(settingRepo *repository.SettingRepository) *SettingService {
	return &SettingService{settingRepo: settingRepo}
}

// GetAllSettings 获取所有设置
// @return 设置列表, 错误
func (s *SettingService) GetAllSettings() (map[string]interface{}, error) {
	settings, err := s.settingRepo.GetAllSettings()
	if err != nil {
		return nil, err
	}

	// 将设置转换为map，并根据类型解析值
	result := make(map[string]interface{})
	for _, setting := range settings {
		parsedValue, err := s.parseValue(setting.Value, setting.Type)
		if err != nil {
			// 如果解析失败，返回原始字符串
			result[setting.Key] = setting.Value
		} else {
			result[setting.Key] = parsedValue
		}
	}

	return result, nil
}

// GetSettingByKey 根据Key获取设置
// @key 设置键
// @return 设置值, 错误
func (s *SettingService) GetSettingByKey(key string) (interface{}, error) {
	setting, err := s.settingRepo.GetSettingByKey(key)
	if err != nil {
		return nil, err
	}

	return s.parseValue(setting.Value, setting.Type)
}

// UpsertSetting 创建或更新设置
// @req 设置请求
// @return 错误
func (s *SettingService) UpsertSetting(req *api.SettingUpsertRequest) error {
	// 验证值是否符合类型
	if err := s.validateValue(req.Value, req.Type); err != nil {
		return err
	}

	setting := &model.Setting{
		Key:   req.Key,
		Value: req.Value,
		Type:  req.Type,
	}

	return s.settingRepo.UpsertSetting(setting)
}

// BatchUpsertSettings 批量创建或更新设置
// @req 批量设置请求
// @return 错误
func (s *SettingService) BatchUpsertSettings(req *api.SettingBatchUpsertRequest) error {
	settings := make([]model.Setting, len(req.Settings))
	for i, settingReq := range req.Settings {
		// 验证值是否符合类型
		if err := s.validateValue(settingReq.Value, settingReq.Type); err != nil {
			return errors.New("设置 " + settingReq.Key + " 的值不符合类型 " + settingReq.Type)
		}

		settings[i] = model.Setting{
			Key:   settingReq.Key,
			Value: settingReq.Value,
			Type:  settingReq.Type,
		}
	}

	return s.settingRepo.BatchUpsertSettings(settings)
}

// DeleteSetting 删除设置
// @key 设置键
// @return 错误
func (s *SettingService) DeleteSetting(key string) error {
	return s.settingRepo.DeleteSetting(key)
}

// parseValue 根据类型解析值
// @value 值
// @valueType 类型
// @return 解析后的值, 错误
func (s *SettingService) parseValue(value string, valueType string) (interface{}, error) {
	switch valueType {
	case "str":
		return value, nil
	case "number":
		// 尝试解析为float64
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, errors.New("无法将值解析为number类型")
		}
		return num, nil
	case "json":
		var result interface{}
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return nil, errors.New("无法将值解析为json类型")
		}
		return result, nil
	default:
		return value, nil
	}
}

// validateValue 验证值是否符合类型
// @value 值
// @valueType 类型
// @return 错误
func (s *SettingService) validateValue(value string, valueType string) error {
	switch valueType {
	case "str":
		return nil
	case "number":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("值不是有效的number类型")
		}
		return nil
	case "json":
		var result interface{}
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return errors.New("值不是有效的json类型")
		}
		return nil
	default:
		return errors.New("不支持的类型: " + valueType)
	}
}
