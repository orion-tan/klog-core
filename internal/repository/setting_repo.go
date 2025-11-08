package repository

import (
	"klog-backend/internal/model"

	"gorm.io/gorm"
)

type SettingRepository struct {
	DB *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{DB: db}
}

// GetSettingByKey 根据Key获取设置
// @key 设置键
// @return 设置, 错误
func (r *SettingRepository) GetSettingByKey(key string) (*model.Setting, error) {
	var setting model.Setting
	err := r.DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetAllSettings 获取所有设置
// @return 设置列表, 错误
func (r *SettingRepository) GetAllSettings() ([]model.Setting, error) {
	var settings []model.Setting
	err := r.DB.Find(&settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpsertSetting 创建或更新设置（如果key存在则更新，不存在则创建）
// @setting 设置
// @return 错误
func (r *SettingRepository) UpsertSetting(setting *model.Setting) error {
	return r.DB.Save(setting).Error
}

// DeleteSetting 删除设置
// @key 设置键
// @return 错误
func (r *SettingRepository) DeleteSetting(key string) error {
	return r.DB.Where("key = ?", key).Delete(&model.Setting{}).Error
}

// BatchUpsertSettings 批量创建或更新设置
// @settings 设置列表
// @return 错误
func (r *SettingRepository) BatchUpsertSettings(settings []model.Setting) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			if err := tx.Save(&setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
