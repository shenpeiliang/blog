package main

import (
	"context"
	"fmt"
	"search/config"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

type Book struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Authors []string  `json:"authors"`
	Price   float64   `json:"price"`
	Genres  []string  `json:"genres"`
	Publish time.Time `json:"publish_date"`
}

type MeiliSearchService struct {
	client  meilisearch.ServiceManager
	indices map[string]meilisearch.IndexManager
}

// 创建新的Meilisearch服务实例
func NewMeiliSearchService(cfg *config.Config) (*MeiliSearchService, error) {
	client := meilisearch.New(cfg.MeiliSearch.Host, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))

	// key 权限校验
	if err := verifyAPIPermissions(client); err != nil {
		return nil, err
	}

	svc := &MeiliSearchService{
		client:  client,
		indices: make(map[string]meilisearch.IndexManager),
	}

	// 初始化所有配置的索引
	for indexName, indexCfg := range cfg.MeiliSearch.Indices {
		// 创建索引
		index := client.Index(indexName)

		// 配置索引设置
		if err := setupIndex(index, &indexCfg); err != nil {
			return nil, fmt.Errorf("初始化索引[%s]失败: %w", indexName, err)
		}
		svc.indices[indexName] = index
	}

	return svc, nil

}

// 验证API权限
func verifyAPIPermissions(client meilisearch.ServiceManager) error {
	// 尝试执行一个低权限操作（例如获取版本信息）
	_, err := client.Version()
	if err != nil {
		return fmt.Errorf("API Key 权限校验失败: %v", err)
	}

	return nil
}

// 配置索引设置
func setupIndex(index meilisearch.IndexManager, cfg *config.IndexConfig) error {
	// 设置可搜索字段
	if _, err := index.UpdateSearchableAttributes(&cfg.SearchableAttributes); err != nil {
		return fmt.Errorf("设置可搜索字段失败: %w", err)
	}

	// 设置可过滤字段
	if _, err := index.UpdateFilterableAttributes(&cfg.FilterableAttributes); err != nil {
		return fmt.Errorf("设置可过滤字段失败: %w", err)
	}

	// 设置排序规则
	if _, err := index.UpdateSortableAttributes(&cfg.SortableAttributes); err != nil {
		return fmt.Errorf("设置排序规则失败: %w", err)
	}

	return nil
}

// 获取索引（安全访问）
func (s *MeiliSearchService) GetIndex(indexName string) (meilisearch.IndexManager, error) {
	if idx, ok := s.indices[indexName]; ok {
		return idx, nil
	}
	return nil, fmt.Errorf("索引[%s]未配置", indexName)
}

// 添加文档到索引
func (s *MeiliSearchService) AddDocuments(indexName string, documents any) (int64, error) {
	//获取索引
	index, err := s.GetIndex(indexName)
	if err != nil {
		return 0, err
	}

	//添加文档
	task, err := index.AddDocuments(documents)
	if err != nil {
		return 0, fmt.Errorf("添加文档到索引失败: %w", err)
	}
	return task.TaskUID, nil
}

// 执行搜索
func (s *MeiliSearchService) Search(indexName string, query string, filter string, sort []string, limit int) ([]any, error) {
	//获取索引
	index, err := s.GetIndex(indexName)
	if err != nil {
		return nil, err
	}

	//执行搜索
	res, err := index.Search(query, &meilisearch.SearchRequest{
		Filter: filter,
		Sort:   sort,
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("搜索失败: %w", err)
	}

	return res.Hits, nil
}

// 等待任务完成
func (s *MeiliSearchService) WaitForTask(taskUID int64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("任务超时")
		default:
			task, err := s.client.GetTask(taskUID)
			if err != nil {
				return fmt.Errorf("获取任务状态失败: %w", err)
			}

			switch task.Status {
			case meilisearch.TaskStatusSucceeded:
				return nil
			case meilisearch.TaskStatusFailed:
				return fmt.Errorf("任务执行失败: %v", task.Error)
			}

			//等待500毫秒
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 删除文档
func (s *MeiliSearchService) DeleteDocument(indexName string, documentID string) (int64, error) {
	//获取索引
	index, err := s.GetIndex(indexName)
	if err != nil {
		return 0, err
	}

	//删除文档
	task, err := index.DeleteDocument(documentID)
	if err != nil {
		return 0, fmt.Errorf("删除文档失败: %w", err)
	}
	return task.TaskUID, nil
}
