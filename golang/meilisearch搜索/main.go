package main

import (
	"encoding/json"
	"fmt"
	"log"
	"search/config"
	"time"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化服务
	service, err := NewMeiliSearchService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	// 添加文档
	books := []Book{
		{
			ID:      "1",
			Title:   "Go语言高级编程",
			Authors: []string{"张三", "李四"},
			Price:   59.9,
			Genres:  []string{"编程", "计算机"},
			Publish: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:      "2",
			Title:   "Rust实战指南",
			Authors: []string{"王五"},
			Price:   69.9,
			Genres:  []string{"编程", "系统"},
			Publish: time.Date(2021, 6, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	// 添加文档
	taskID, err := service.AddDocuments("books", books)
	if err != nil {
		log.Fatalf("Failed to add documents: %v", err)
	}

	// 等待任务完成方便后续搜索
	if err := service.WaitForTask(taskID, 10*time.Second); err != nil {
		log.Fatalf("Task failed: %v", err)
	}

	// 执行搜索
	results, err := service.Search("books", "编程", "price < 70", []string{"price:asc"}, 10)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	var booksResults []Book
	for _, hit := range results {
		data, err := json.Marshal(hit)
		if err != nil {
			continue // 跳过无法解析的结果
		}

		var book Book
		if err := json.Unmarshal(data, &book); err == nil {
			booksResults = append(booksResults, book)
		}
	}

	// 处理结果
	fmt.Printf("Found %d results:\n", len(results))
	for _, book := range booksResults {
		fmt.Printf("ID: %s, Title: %s, Price: %.2f\n", book.ID, book.Title, book.Price)
	}

	// 删除文档示例
	// if _, err := service.DeleteDocument("1"); err != nil {
	// 	log.Printf("Failed to delete document: %v", err)
	// }
}
