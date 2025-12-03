package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/XiaoLFeng/llm-memory/internal/models"
	"github.com/XiaoLFeng/llm-memory/internal/models/dto"
	"github.com/XiaoLFeng/llm-memory/internal/models/entity"
	"github.com/XiaoLFeng/llm-memory/pkg/types"
)

// ToDoService 待办事项服务
// 嘿嘿~ 这是处理待办事项业务逻辑的服务层哦！💖
// 注意：类型名使用 ToDo，MCP 工具名保持 todo_*
type ToDoService struct {
	model *models.ToDoModel
}

// NewToDoService 创建新的待办事项服务实例
// 呀~ 构造函数来啦！(´∀｀)
func NewToDoService(model *models.ToDoModel) *ToDoService {
	return &ToDoService{
		model: model,
	}
}

// CreateToDo 创建新的待办事项
// 嘿嘿~ 创建待办前会先验证数据的完整性呢！💫
func (s *ToDoService) CreateToDo(ctx context.Context, input *dto.ToDoCreateDTO, scopeCtx *types.ScopeContext) (*entity.ToDo, error) {
	// 验证标题不能为空
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("标题不能为空哦~ 📝")
	}

	// 默认优先级
	priority := entity.ToDoPriority(input.Priority)
	if priority < entity.ToDoPriorityLow || priority > entity.ToDoPriorityUrgent {
		priority = entity.ToDoPriorityMedium
	}

	// 解析作用域
	var groupID uint
	var path string

	scope := strings.ToLower(input.Scope)
	switch scope {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
	case "global":
		// groupID 和 path 都为空即为 global
	default:
		// 默认：group 优先，然后 personal
		groupID, path = resolveDefaultScope(scopeCtx)
	}

	// 创建待办事项实例
	todo := &entity.ToDo{
		GroupID:     groupID,
		Path:        path,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Priority:    priority,
		Status:      entity.ToDoStatusPending,
		DueDate:     input.DueDate,
	}

	// 保存到数据库
	if err := s.model.Create(ctx, todo); err != nil {
		return nil, err
	}

	// 更新标签
	if len(input.Tags) > 0 {
		if err := s.model.UpdateTags(ctx, todo.ID, input.Tags); err != nil {
			return nil, err
		}
		// 重新获取以包含标签
		todo, _ = s.model.FindByID(ctx, todo.ID)
	}

	return todo, nil
}

// UpdateToDo 更新待办事项
func (s *ToDoService) UpdateToDo(ctx context.Context, input *dto.ToDoUpdateDTO) error {
	// 验证ID
	if input.ID == 0 {
		return errors.New("待办事项ID不能为0")
	}

	// 获取现有待办
	todo, err := s.model.FindByID(ctx, input.ID)
	if err != nil {
		return errors.New("待办事项不存在")
	}

	// 应用更新
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return errors.New("标题不能为空")
		}
		todo.Title = title
	}
	if input.Description != nil {
		todo.Description = strings.TrimSpace(*input.Description)
	}
	if input.Priority != nil {
		priority := entity.ToDoPriority(*input.Priority)
		if priority < entity.ToDoPriorityLow || priority > entity.ToDoPriorityUrgent {
			return errors.New("无效的优先级")
		}
		todo.Priority = priority
	}
	if input.Status != nil {
		status := entity.ToDoStatus(*input.Status)
		todo.Status = status
		if status == entity.ToDoStatusCompleted && todo.CompletedAt == nil {
			now := time.Now()
			todo.CompletedAt = &now
		}
	}
	if input.DueDate != nil {
		todo.DueDate = input.DueDate
	}

	// 执行更新
	if err := s.model.Update(ctx, todo); err != nil {
		return err
	}

	// 更新标签（如果提供）
	if input.Tags != nil {
		if err := s.model.UpdateTags(ctx, todo.ID, *input.Tags); err != nil {
			return err
		}
	}

	return nil
}

// DeleteToDo 删除待办事项
func (s *ToDoService) DeleteToDo(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	// 检查是否存在
	_, err := s.model.FindByID(ctx, id)
	if err != nil {
		return errors.New("待办事项不存在")
	}

	return s.model.Delete(ctx, id)
}

// GetToDo 获取指定ID的待办事项
func (s *ToDoService) GetToDo(ctx context.Context, id uint) (*entity.ToDo, error) {
	if id == 0 {
		return nil, errors.New("无效的待办事项ID")
	}

	return s.model.FindByID(ctx, id)
}

// ListToDos 获取所有待办事项
func (s *ToDoService) ListToDos(ctx context.Context) ([]entity.ToDo, error) {
	return s.model.FindAll(ctx)
}

// ListToDosByScope 根据作用域列出待办事项
// 嘿嘿~ 支持 Personal/Group/Global 三层作用域过滤！💖
func (s *ToDoService) ListToDosByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.ToDo, error) {
	var groupID uint
	var path string
	var includeGlobal bool

	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
		includeGlobal = false
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
		includeGlobal = false
	case "global":
		includeGlobal = true
	case "all", "":
		if scopeCtx != nil {
			if scopeCtx.CurrentPath != "" {
				path = scopeCtx.CurrentPath
			}
			if scopeCtx.GroupID > 0 {
				groupID = uint(scopeCtx.GroupID)
			}
		}
		includeGlobal = true
	default:
		includeGlobal = true
	}

	return s.model.FindByScope(ctx, groupID, path, includeGlobal)
}

// ListByStatus 根据状态获取待办事项列表
func (s *ToDoService) ListByStatus(ctx context.Context, status entity.ToDoStatus) ([]entity.ToDo, error) {
	return s.model.FindByStatus(ctx, status)
}

// ListToday 获取今天的待办事项
func (s *ToDoService) ListToday(ctx context.Context) ([]entity.ToDo, error) {
	return s.model.FindToday(ctx)
}

// ListTodayByScope 根据作用域获取今天的待办事项
func (s *ToDoService) ListTodayByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.ToDo, error) {
	var groupID uint
	var path string
	var includeGlobal bool

	switch strings.ToLower(scope) {
	case "personal":
		if scopeCtx != nil && scopeCtx.CurrentPath != "" {
			path = scopeCtx.CurrentPath
		}
		includeGlobal = false
	case "group":
		if scopeCtx != nil && scopeCtx.GroupID > 0 {
			groupID = uint(scopeCtx.GroupID)
		}
		includeGlobal = false
	case "global":
		includeGlobal = true
	case "all", "":
		if scopeCtx != nil {
			if scopeCtx.CurrentPath != "" {
				path = scopeCtx.CurrentPath
			}
			if scopeCtx.GroupID > 0 {
				groupID = uint(scopeCtx.GroupID)
			}
		}
		includeGlobal = true
	default:
		includeGlobal = true
	}

	return s.model.FindTodayByScope(ctx, groupID, path, includeGlobal)
}

// CompleteToDo 标记待办事项为已完成
func (s *ToDoService) CompleteToDo(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("待办事项已经完成")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法完成")
	}

	return s.model.Complete(ctx, id)
}

// StartToDo 标记待办事项为进行中
func (s *ToDoService) StartToDo(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法开始")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法开始")
	}

	return s.model.Start(ctx, id)
}

// CancelToDo 取消待办事项
func (s *ToDoService) CancelToDo(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.model.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法取消")
	}

	return s.model.Cancel(ctx, id)
}

// ========== 批量操作方法 ==========

// BatchCreateToDos 批量创建待办事项
// 嘿嘿~ 一次性创建多个待办！🎮
func (s *ToDoService) BatchCreateToDos(ctx context.Context, input *dto.ToDoBatchCreateDTO, scopeCtx *types.ScopeContext) (*dto.ToDoBatchResultDTO, error) {
	// 验证数量限制
	if len(input.Items) == 0 {
		return nil, errors.New("没有待创建的项目")
	}
	if len(input.Items) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	// 转换为 entity 列表
	todos := make([]entity.ToDo, 0, len(input.Items))
	for _, item := range input.Items {
		if strings.TrimSpace(item.Title) == "" {
			continue // 跳过空标题
		}

		// 解析作用域
		var groupID uint
		var path string

		scope := strings.ToLower(item.Scope)
		switch scope {
		case "personal":
			if scopeCtx != nil && scopeCtx.CurrentPath != "" {
				path = scopeCtx.CurrentPath
			}
		case "group":
			if scopeCtx != nil && scopeCtx.GroupID > 0 {
				groupID = uint(scopeCtx.GroupID)
			}
		case "global":
			// 留空
		default:
			groupID, path = resolveDefaultScope(scopeCtx)
		}

		priority := entity.ToDoPriority(item.Priority)
		if priority < entity.ToDoPriorityLow || priority > entity.ToDoPriorityUrgent {
			priority = entity.ToDoPriorityMedium
		}

		todo := entity.ToDo{
			GroupID:     groupID,
			Path:        path,
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(item.Description),
			Priority:    priority,
			Status:      entity.ToDoStatusPending,
			DueDate:     item.DueDate,
		}
		todos = append(todos, todo)
	}

	if len(todos) == 0 {
		return nil, errors.New("没有有效的待创建项目")
	}

	return s.model.BatchCreate(ctx, todos)
}

// BatchUpdateToDos 批量更新待办事项
// 呀~ 一次性更新多个待办！✨
func (s *ToDoService) BatchUpdateToDos(ctx context.Context, input *dto.ToDoBatchUpdateDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("没有待更新的项目")
	}
	if len(input.Items) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.model.BatchUpdate(ctx, input.Items)
}

// BatchCompleteToDos 批量完成待办事项
// 嘿嘿~ 一次性完成多个待办！💖
func (s *ToDoService) BatchCompleteToDos(ctx context.Context, input *dto.ToDoBatchCompleteDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.IDs) == 0 {
		return nil, errors.New("没有待完成的项目")
	}
	if len(input.IDs) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.model.BatchComplete(ctx, input.IDs)
}

// BatchDeleteToDos 批量删除待办事项
// 呀~ 一次性删除多个待办！⚠️
func (s *ToDoService) BatchDeleteToDos(ctx context.Context, input *dto.ToDoBatchDeleteDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.IDs) == 0 {
		return nil, errors.New("没有待删除的项目")
	}
	if len(input.IDs) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.model.BatchDelete(ctx, input.IDs)
}

// BatchUpdateStatus 批量更新状态
func (s *ToDoService) BatchUpdateStatus(ctx context.Context, ids []uint, status entity.ToDoStatus) (*dto.ToDoBatchResultDTO, error) {
	if len(ids) == 0 {
		return nil, errors.New("没有待更新的项目")
	}
	if len(ids) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.model.BatchUpdateStatus(ctx, ids, status)
}

// ToToDoResponseDTO 将 ToDo entity 转换为 ResponseDTO
// 嘿嘿~ 数据转换小助手！💖
func ToToDoResponseDTO(todo *entity.ToDo, currentPath string) *dto.ToDoResponseDTO {
	if todo == nil {
		return nil
	}

	tags := make([]string, 0, len(todo.Tags))
	for _, t := range todo.Tags {
		tags = append(tags, t.Tag)
	}

	// 判断作用域
	var scope types.Scope
	if todo.Path != "" {
		scope = types.ScopePersonal
	} else if todo.GroupID > 0 {
		scope = types.ScopeGroup
	} else {
		scope = types.ScopeGlobal
	}

	return &dto.ToDoResponseDTO{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    int(todo.Priority),
		Status:      int(todo.Status),
		DueDate:     todo.DueDate,
		CompletedAt: todo.CompletedAt,
		Tags:        tags,
		Scope:       string(scope),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

// ToToDoListDTO 将 ToDo entity 转换为 ListDTO
func ToToDoListDTO(todo *entity.ToDo) *dto.ToDoListDTO {
	if todo == nil {
		return nil
	}

	return &dto.ToDoListDTO{
		ID:       todo.ID,
		Title:    todo.Title,
		Priority: int(todo.Priority),
		Status:   int(todo.Status),
		DueDate:  todo.DueDate,
	}
}
