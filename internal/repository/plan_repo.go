package repository

import (
	"context"

	"github.com/XiaoLFeng/llm-memory/internal/database"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// PlanRepo 计划仓储实现结构体
// 嘿嘿~ 这是计划管理的核心实现呢！📋
type PlanRepo struct {
	dbPath string
}

// NewPlanRepo 创建新的计划仓储实例
// 构造函数模式，让代码更优雅~ 💖
func NewPlanRepo(dbPath string) *PlanRepo {
	return &PlanRepo{
		dbPath: dbPath,
	}
}

// Create 创建新的计划
// 使用 db.Save 保存计划到数据库~ ✨
func (r *PlanRepo) Create(ctx context.Context, plan *types.Plan) error {
	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		return db.Save(plan)
	})
}

// Update 更新现有计划
// 使用 db.Update 更新计划信息~ 🎮
func (r *PlanRepo) Update(ctx context.Context, plan *types.Plan) error {
	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		return db.Update(plan)
	})
}

// Delete 删除指定ID的计划
// 使用 db.DeleteStruct 删除计划~ (´∀｀)
func (r *PlanRepo) Delete(ctx context.Context, id int) error {
	return database.OpenWithActionNoReturn(r.dbPath, func(db *database.DB) error {
		return db.DeleteStruct(&types.Plan{ID: id})
	})
}

// FindByID 根据ID查找计划
// 使用 db.One 查询单个计划~ 🍫
func (r *PlanRepo) FindByID(ctx context.Context, id int) (*types.Plan, error) {
	return database.OpenWithAction(r.dbPath, func(db *database.DB) (*types.Plan, error) {
		var plan types.Plan
		err := db.One("ID", id, &plan)
		if err != nil {
			return nil, err
		}
		return &plan, nil
	})
}

// FindAll 查找所有计划
// 使用 db.All 获取全部计划列表~ ＼(^o^)／
func (r *PlanRepo) FindAll(ctx context.Context) ([]types.Plan, error) {
	return database.OpenWithAction(r.dbPath, func(db *database.DB) ([]types.Plan, error) {
		var plans []types.Plan
		err := db.All(&plans)
		if err != nil {
			return nil, err
		}
		return plans, nil
	})
}

// FindByStatus 根据状态查找计划
// 使用 db.Find 按状态筛选计划~ 🎯
func (r *PlanRepo) FindByStatus(ctx context.Context, status types.PlanStatus) ([]types.Plan, error) {
	return database.OpenWithAction(r.dbPath, func(db *database.DB) ([]types.Plan, error) {
		var plans []types.Plan
		err := db.Find("Status", status, &plans)
		if err != nil {
			return nil, err
		}
		return plans, nil
	})
}

// FindByScope 根据作用域查找计划
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
func (r *PlanRepo) FindByScope(ctx context.Context, scope *types.ScopeContext) ([]types.Plan, error) {
	if scope == nil {
		// 没有作用域限制，返回所有
		return r.FindAll(ctx)
	}

	return database.OpenWithAction(r.dbPath, func(db *database.DB) ([]types.Plan, error) {
		var allPlans []types.Plan
		err := db.All(&allPlans)
		if err != nil {
			return nil, err
		}

		var result []types.Plan
		for _, plan := range allPlans {
			if r.matchScope(plan, scope) {
				result = append(result, plan)
			}
		}

		return result, nil
	})
}

// matchScope 检查计划是否匹配作用域
// 核心过滤逻辑~ ✨
func (r *PlanRepo) matchScope(plan types.Plan, scope *types.ScopeContext) bool {
	// 检查 Global
	if scope.IncludeGlobal && plan.IsGlobal() {
		return true
	}

	// 检查 Personal（精确路径匹配）
	if scope.IncludePersonal && plan.Path != "" && plan.Path == scope.CurrentPath {
		return true
	}

	// 检查 Group
	if scope.IncludeGroup && scope.GroupID != types.GlobalGroupID && plan.GroupID == scope.GroupID {
		return true
	}

	return false
}
