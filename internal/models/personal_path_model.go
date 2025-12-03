package models

import (
	"context"
	"path/filepath"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"gorm.io/gorm"
)

// PersonalPathModel 个人路径数据访问层
type PersonalPathModel struct {
	db *gorm.DB
}

// NewPersonalPathModel 创建 PersonalPathModel 实例
func NewPersonalPathModel(db *gorm.DB) *PersonalPathModel {
	return &PersonalPathModel{db: db}
}

// EnsurePath 确保路径存在，不存在则创建，存在则更新访问时间
func (m *PersonalPathModel) EnsurePath(ctx context.Context, path string) (*entity.PersonalPath, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		personalPath = entity.PersonalPath{
			ID:        database.GenerateID(),
			Path:      absPath,
			LastVisit: time.Now(),
		}
		if err := m.db.WithContext(ctx).Create(&personalPath).Error; err != nil {
			return nil, err
		}
		return &personalPath, nil
	}

	if err != nil {
		return nil, err
	}

	// 已存在，更新访问时间
	personalPath.Touch()
	if err := m.db.WithContext(ctx).Save(&personalPath).Error; err != nil {
		return nil, err
	}
	return &personalPath, nil
}

// FindByPath 根据路径查找记录
func (m *PersonalPathModel) FindByPath(ctx context.Context, path string) (*entity.PersonalPath, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var personalPath entity.PersonalPath
	err = m.db.WithContext(ctx).Where("path = ?", absPath).First(&personalPath).Error
	if err != nil {
		return nil, err
	}
	return &personalPath, nil
}

// Exists 检查路径是否已存在
func (m *PersonalPathModel) Exists(ctx context.Context, path string) (bool, error) {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var count int64
	err = m.db.WithContext(ctx).Model(&entity.PersonalPath{}).Where("path = ?", absPath).Count(&count).Error
	return count > 0, err
}

// ListRecentPaths 列出最近访问的路径
func (m *PersonalPathModel) ListRecentPaths(ctx context.Context, limit int) ([]entity.PersonalPath, error) {
	if limit <= 0 {
		limit = 10
	}

	var paths []entity.PersonalPath
	err := m.db.WithContext(ctx).Order("last_visit DESC").Limit(limit).Find(&paths).Error
	return paths, err
}

// Delete 删除路径记录（软删除）
func (m *PersonalPathModel) Delete(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Delete(&entity.PersonalPath{}, id).Error
}

// DeleteByPath 根据路径删除记录
func (m *PersonalPathModel) DeleteByPath(ctx context.Context, path string) error {
	// 规范化路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return m.db.WithContext(ctx).Where("path = ?", absPath).Delete(&entity.PersonalPath{}).Error
}

// Count 获取路径总数
func (m *PersonalPathModel) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.db.WithContext(ctx).Model(&entity.PersonalPath{}).Count(&count).Error
	return count, err
}
