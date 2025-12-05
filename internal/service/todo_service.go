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
// 注意：类型名使用 ToDo，MCP 工具名保持 todo_*
type ToDoService struct {
	todoModel *models.ToDoModel
}

// NewToDoService 创建新的待办事项服务实例
func NewToDoService(model *models.ToDoModel) *ToDoService {
	return &ToDoService{
		todoModel: model,
	}
}

// CreateToDo 创建新的待办事项
// PathID 关联个人或小组路径
func (s *ToDoService) CreateToDo(ctx context.Context, input *dto.ToDoCreateDTO, scopeCtx *types.ScopeContext) (*entity.ToDo, error) {
	// 验证标题不能为空
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("标题不能为空")
	}

	// 验证 code 格式
	if err := entity.ValidateCode(input.Code); err != nil {
		return nil, err
	}

	// 检查活跃状态中 code 唯一性
	exists, err := s.todoModel.ExistsActiveCode(ctx, input.Code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("代码已被活跃待办使用")
	}

	// 默认优先级
	priority := entity.ToDoPriority(input.Priority)
	if priority < entity.ToDoPriorityLow || priority > entity.ToDoPriorityUrgent {
		priority = entity.ToDoPriorityMedium
	}

	// 解析作用域 -> PathID（必须指定路径）
	pathID := resolveDefaultPathID(scopeCtx)
	if pathID == 0 {
		return nil, errors.New("无法确定作用域，请先初始化 paths")
	}

	// 创建待办事项实例
	todo := &entity.ToDo{
		PathID:      pathID,
		Code:        input.Code,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Priority:    priority,
		Status:      entity.ToDoStatusPending,
		DueDate:     input.DueDate,
	}

	// 保存到数据库
	if err := s.todoModel.Create(ctx, todo); err != nil {
		return nil, err
	}

	// 更新标签
	if len(input.Tags) > 0 {
		if err := s.todoModel.UpdateTags(ctx, todo.ID, input.Tags); err != nil {
			return nil, err
		}
		// 重新获取以包含标签
		todo, _ = s.todoModel.FindByID(ctx, todo.ID)
	}

	return todo, nil
}

// UpdateToDo 更新待办事项
func (s *ToDoService) UpdateToDo(ctx context.Context, input *dto.ToDoUpdateDTO) error {
	// 验证 Code
	if input.Code == "" {
		return errors.New("待办事项 Code 不能为空")
	}

	// 通过 Code 获取现有待办
	todo, err := s.todoModel.FindByCode(ctx, input.Code)
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
	if err := s.todoModel.Update(ctx, todo); err != nil {
		return err
	}

	// 更新标签（如果提供）
	if input.Tags != nil {
		if err := s.todoModel.UpdateTags(ctx, todo.ID, *input.Tags); err != nil {
			return err
		}
	}

	return nil
}

// DeleteToDo 删除待办事项
func (s *ToDoService) DeleteToDo(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("无效的待办事项 Code")
	}

	// 通过 Code 获取待办
	todo, err := s.todoModel.FindByCode(ctx, code)
	if err != nil {
		return errors.New("待办事项不存在")
	}

	return s.todoModel.Delete(ctx, todo.ID)
}

// DeleteToDoByID 根据 ID 删除待办（TUI 内部使用）
func (s *ToDoService) DeleteToDoByID(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	// 检查是否存在
	_, err := s.todoModel.FindByID(ctx, id)
	if err != nil {
		return errors.New("待办事项不存在")
	}

	return s.todoModel.Delete(ctx, id)
}

// GetToDo 获取指定 Code 的待办事项
func (s *ToDoService) GetToDo(ctx context.Context, code string) (*entity.ToDo, error) {
	if code == "" {
		return nil, errors.New("无效的待办事项 Code")
	}

	return s.todoModel.FindByCode(ctx, code)
}

// GetToDoByID 根据 ID 获取待办（TUI 内部使用）
func (s *ToDoService) GetToDoByID(ctx context.Context, id int64) (*entity.ToDo, error) {
	if id == 0 {
		return nil, errors.New("无效的待办事项ID")
	}

	return s.todoModel.FindByID(ctx, id)
}

// ListToDos 获取所有待办事项（需要提供作用域上下文）
// 使用 PathOnlyFilter 进行过滤
func (s *ToDoService) ListToDos(ctx context.Context, scopeCtx *types.ScopeContext) ([]entity.ToDo, error) {
	filter := buildPathOnlyFilter("all", scopeCtx)
	return s.todoModel.FindByPathOnlyFilter(ctx, filter)
}

// ListToDosByScope 根据作用域列出待办事项
// 使用 PathOnlyFilter 进行过滤（无 Global 支持）
func (s *ToDoService) ListToDosByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) ([]entity.ToDo, error) {
	filter := buildPathOnlyFilter(scope, scopeCtx)
	return s.todoModel.FindByPathOnlyFilter(ctx, filter)
}

// ListByStatus 根据状态获取待办事项列表
func (s *ToDoService) ListByStatus(ctx context.Context, status entity.ToDoStatus, scopeCtx *types.ScopeContext) ([]entity.ToDo, error) {
	filter := buildPathOnlyFilter("all", scopeCtx)
	return s.todoModel.FindByStatus(ctx, status, filter)
}

// CompleteToDo 标记待办事项为已完成
func (s *ToDoService) CompleteToDo(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("无效的待办事项 Code")
	}

	todo, err := s.todoModel.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("待办事项已经完成")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法完成")
	}

	return s.todoModel.Complete(ctx, todo.ID)
}

// CompleteToDoByID 根据 ID 完成待办（TUI 内部使用）
func (s *ToDoService) CompleteToDoByID(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.todoModel.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("待办事项已经完成")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法完成")
	}

	return s.todoModel.Complete(ctx, id)
}

// StartToDo 标记待办事项为进行中
func (s *ToDoService) StartToDo(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("无效的待办事项 Code")
	}

	todo, err := s.todoModel.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法开始")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法开始")
	}

	return s.todoModel.Start(ctx, todo.ID)
}

// StartToDoByID 根据 ID 开始待办（TUI 内部使用）
func (s *ToDoService) StartToDoByID(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.todoModel.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法开始")
	}
	if todo.Status == entity.ToDoStatusCancelled {
		return errors.New("已取消的待办事项无法开始")
	}

	return s.todoModel.Start(ctx, id)
}

// CancelToDo 取消待办事项
func (s *ToDoService) CancelToDo(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("无效的待办事项 Code")
	}

	todo, err := s.todoModel.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法取消")
	}

	return s.todoModel.Cancel(ctx, todo.ID)
}

// CancelToDoByID 根据 ID 取消待办（TUI 内部使用）
func (s *ToDoService) CancelToDoByID(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("无效的待办事项ID")
	}

	todo, err := s.todoModel.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.Status == entity.ToDoStatusCompleted {
		return errors.New("已完成的待办事项无法取消")
	}

	return s.todoModel.Cancel(ctx, id)
}

// BatchCreateToDos 批量创建待办事项
// PathID 关联个人或小组路径
func (s *ToDoService) BatchCreateToDos(ctx context.Context, input *dto.ToDoBatchCreateDTO, scopeCtx *types.ScopeContext) (*dto.ToDoBatchResultDTO, error) {
	// 验证数量限制
	if len(input.Items) == 0 {
		return nil, errors.New("没有待创建的项目")
	}
	if len(input.Items) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	// 解析作用域 -> PathID（必须指定路径）
	pathID := resolveDefaultPathID(scopeCtx)
	if pathID == 0 {
		return nil, errors.New("无法确定作用域，请先初始化 paths")
	}

	// 转换为 entity 列表
	todos := make([]entity.ToDo, 0, len(input.Items))
	for _, item := range input.Items {
		if strings.TrimSpace(item.Title) == "" {
			continue // 跳过空标题
		}

		priority := entity.ToDoPriority(item.Priority)
		if priority < entity.ToDoPriorityLow || priority > entity.ToDoPriorityUrgent {
			priority = entity.ToDoPriorityMedium
		}

		todo := entity.ToDo{
			PathID:      pathID,
			Code:        item.Code,
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

	return s.todoModel.BatchCreate(ctx, todos)
}

// BatchUpdateToDos 批量更新待办事项
func (s *ToDoService) BatchUpdateToDos(ctx context.Context, input *dto.ToDoBatchUpdateDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("没有待更新的项目")
	}
	if len(input.Items) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.todoModel.BatchUpdate(ctx, input.Items)
}

// BatchCompleteToDos 批量完成待办事项
func (s *ToDoService) BatchCompleteToDos(ctx context.Context, input *dto.ToDoBatchCompleteDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.Codes) == 0 {
		return nil, errors.New("没有待完成的项目")
	}
	if len(input.Codes) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	// 将 Codes 转换为 IDs
	ids := make([]int64, 0, len(input.Codes))
	for _, code := range input.Codes {
		todo, err := s.todoModel.FindByCode(ctx, code)
		if err == nil && todo != nil {
			ids = append(ids, todo.ID)
		}
	}

	if len(ids) == 0 {
		return nil, errors.New("未找到有效的待办事项")
	}

	return s.todoModel.BatchComplete(ctx, ids)
}

// BatchDeleteToDos 批量删除待办事项
func (s *ToDoService) BatchDeleteToDos(ctx context.Context, input *dto.ToDoBatchDeleteDTO) (*dto.ToDoBatchResultDTO, error) {
	if len(input.Codes) == 0 {
		return nil, errors.New("没有待删除的项目")
	}
	if len(input.Codes) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	// 将 Codes 转换为 IDs
	ids := make([]int64, 0, len(input.Codes))
	for _, code := range input.Codes {
		todo, err := s.todoModel.FindByCode(ctx, code)
		if err == nil && todo != nil {
			ids = append(ids, todo.ID)
		}
	}

	if len(ids) == 0 {
		return nil, errors.New("未找到有效的待办事项")
	}

	return s.todoModel.BatchDelete(ctx, ids)
}

// BatchUpdateStatus 批量更新状态
func (s *ToDoService) BatchUpdateStatus(ctx context.Context, ids []int64, status entity.ToDoStatus) (*dto.ToDoBatchResultDTO, error) {
	if len(ids) == 0 {
		return nil, errors.New("没有待更新的项目")
	}
	if len(ids) > dto.MaxBatchSize {
		return nil, errors.New("批量操作最多支持 100 条记录")
	}

	return s.todoModel.BatchUpdateStatus(ctx, ids, status)
}

// DeleteAllByScope 删除当前作用域内的所有待办事项（用于 todo_final）
// 返回删除的记录数量
func (s *ToDoService) DeleteAllByScope(ctx context.Context, scope string, scopeCtx *types.ScopeContext) (int64, error) {
	filter := buildPathOnlyFilter(scope, scopeCtx)
	return s.todoModel.BatchDeleteByPathIDs(ctx, filter.PathIDs)
}

// ToToDoResponseDTO 将 ToDo entity 转换为 ResponseDTO
// 使用 PathID 判断作用域（无 Global 支持）
func ToToDoResponseDTO(todo *entity.ToDo, scopeCtx *types.ScopeContext) *dto.ToDoResponseDTO {
	if todo == nil {
		return nil
	}

	tags := make([]string, 0, len(todo.Tags))
	for _, t := range todo.Tags {
		tags = append(tags, t.Tag)
	}

	scope := types.GetScopeForDisplayNoGlobal(todo.PathID, scopeCtx)

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
