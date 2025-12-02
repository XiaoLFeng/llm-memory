package repository

import (
	"context"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// PlanRepo 计划仓储实现结构体
// 嘿嘿~ 这是计划管理的核心实现呢！📋
type PlanRepo struct {
	db *database.DB
}

// NewPlanRepo 创建新的计划仓储实例
// 构造函数模式，让代码更优雅~ 💖
func NewPlanRepo(db *database.DB) *PlanRepo {
	return &PlanRepo{
		db: db,
	}
}

// Create 创建新的计划
// 使用 db.Save 保存计划到数据库~ ✨
func (r *PlanRepo) Create(ctx context.Context, plan *types.Plan) error {
	return r.db.Save(plan)
}

// Update 更新现有计划
// 使用 db.Update 更新计划信息~ 🎮
func (r *PlanRepo) Update(ctx context.Context, plan *types.Plan) error {
	return r.db.Update(plan)
}

// Delete 删除指定ID的计划
// 使用 db.DeleteStruct 删除计划~ (´∀｀)
func (r *PlanRepo) Delete(ctx context.Context, id int) error {
	return r.db.DeleteStruct(&types.Plan{ID: id})
}

// FindByID 根据ID查找计划
// 使用 db.One 查询单个计划~ 🍫
func (r *PlanRepo) FindByID(ctx context.Context, id int) (*types.Plan, error) {
	var plan types.Plan
	err := r.db.One("ID", id, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindAll 查找所有计划
// 使用 db.All 获取全部计划列表~ ＼(^o^)／
func (r *PlanRepo) FindAll(ctx context.Context) ([]types.Plan, error) {
	var plans []types.Plan
	err := r.db.All(&plans)
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// FindByStatus 根据状态查找计划
// 使用 db.Find 按状态筛选计划~ 🎯
func (r *PlanRepo) FindByStatus(ctx context.Context, status types.PlanStatus) ([]types.Plan, error) {
	var plans []types.Plan
	err := r.db.Find("Status", status, &plans)
	if err != nil {
		return nil, err
	}
	return plans, nil
}
