# MeiliSearch 使用文档

## 目录
### 基础篇
- [下载与安装](#下载与安装)
- [基本配置](#基本配置)
- [启动服务](#启动服务)
- [API Key管理](#api-key管理)
- [验证服务](#验证服务)

### Go客户端篇
- [配置说明](#配置说明)
  - [基础配置](#基础配置)
- [服务层实现](#服务层实现)
  - [初始化服务](#初始化服务)
  - [数据结构](#数据结构)
- [核心功能](#核心功能)
  - [文档操作](#文档操作)
  - [搜索功能](#搜索功能)
- [使用示例](#使用示例)
  - [完整示例代码](#完整示例代码)

## 下载与安装

从 GitHub 下载最新版本的 MeiliSearch：
- 下载地址：https://github.com/meilisearch/MeiliSearch/releases

## 基本配置

### 关键概念
- **master-key**：最高权限密钥，用于生成其他密钥，通常在生产环境中使用
- **API Key**：普通操作密钥，用于客户端连接

### 默认配置
- 默认监听地址：`http://127.0.0.1:7700`

### 常用启动参数
| 参数 | 描述 | 示例 |
|------|------|------|
| `--env` | 指定环境模式 | `--env production` |
| `--db-path` | 指定数据存储目录 | `--db-path C:\meilisearch\data` |
| `--http-addr` | 指定监听地址 | `--http-addr 0.0.0.0:8080` |
| `--master-key` | 设置主密钥 | `--master-key="your_master_key_here"` |

## 启动服务

### Windows 示例
```bash
./meilisearch-windows-amd64.exe --master-key="pU3GxrEG_QOl7Ky68cLrN8TJjT-wzRF1BjmUWKV_tDY" --env="production"
```

### Linux 示例
```bash
./meilisearch --master-key="your_master_key" --env="production" --db-path="/var/lib/meilisearch/data.ms"
```

## API Key管理

### 生成API Key
```bash
curl -X POST 'http://localhost:7700/keys' \
  -H 'Authorization: Bearer your_master_key_here' \
  -H 'Content-Type: application/json' \
  --data-binary '{
    "name": "图书",
    "description": "图书内容搜索",
    "actions": ["search", "documents.add", "documents.delete", "indexes.create", "indexes.delete", "tasks.get", "settings.*"],
    "indexes": ["*"],
    "expiresAt": null
  }'
```

注意：中文需要Unicode转码

### 列出所有API Key
```bash
curl -X GET 'http://localhost:7700/keys' \
  -H 'Authorization: Bearer your_master_key_here' | \
  sed 's/{/{\n/g' | sed 's/}/}\n/g' | sed 's/,/, /g' | awk '{print "  "$0}'
```

### 删除API Key
```bash
curl -X DELETE 'http://localhost:7700/keys/key_uid_here' \
  -H 'Authorization: Bearer your_master_key_here'
```

### 测试Key权限
```bash
curl -X GET 'http://localhost:7700/version' \
  -H 'Authorization: Bearer your_api_key_here'
```

## 验证服务

打开浏览器访问默认地址：
```
http://localhost:7700
```

## 官方文档

更多详细信息和高级功能请参考官方文档：
https://www.meilisearch.com/docs/learn/filtering_and_sorting/working_with_dates

API Key管理参考：
https://www.meilisearch.com/docs/reference/api/keys


# MeiliSearch Go 客户端使用

## 配置说明

### 基础配置
```toml
[meilisearch]
host = "http://localhost:7700"
master_key = "pU3GxrEG_QOl7Ky68cLrN8TJjT-wzRF1BjmUWKV_tDY" 
api_key = "b6a1552edabd1f649535df47d6923fa3f4b9f4d7b2c4f567caf5e5009aeeeeed"

[meilisearch.indices.books]
searchable_attributes = ["title", "authors", "genres"] # 可搜索字段
filterable_attributes = ["price", "publish_date"] # 可过滤字段
sortable_attributes = ["price", "publish_date"] # 可排序字段
```


## 服务层实现

### 初始化服务
```go
func NewMeiliSearchService(cfg *config.Config) (*MeiliSearchService, error) {
    client := meilisearch.New(cfg.MeiliSearch.Host, meilisearch.WithAPIKey(cfg.MeiliSearch.APIKey))
    
    // 权限校验
    if err := verifyAPIPermissions(client); err != nil {
        return nil, err
    }
    
    // 初始化索引配置
    for indexName, indexCfg := range cfg.MeiliSearch.Indices {
        index := client.Index(indexName)
        if err := setupIndex(index, &indexCfg); err != nil {
            return nil, fmt.Errorf("初始化索引[%s]失败: %w", indexName, err)
        }
    }
    
    return &MeiliSearchService{client: client}, nil
}
```

### 数据结构
```go
type Book struct {
    ID      string    `json:"id"`
    Title   string    `json:"title"`
    Authors []string  `json:"authors"`
    Price   float64   `json:"price"`
    Genres  []string  `json:"genres"`
    Publish time.Time `json:"publish_date"`
}
```

## 核心功能

### 文档操作
```go
// 添加文档
taskID, err := service.AddDocuments("books", books)

// 删除文档
taskID, err := service.DeleteDocument("books", "1")

// 等待任务完成
err := service.WaitForTask(taskID, 10*time.Second)
```

### 搜索功能
```go
// 基本搜索
results, err := service.Search("books", "编程", "", nil, 10)

// 带过滤条件
results, err := service.Search("books", "编程", "price < 70", nil, 10)

// 带排序
results, err := service.Search("books", "编程", "", []string{"price:asc"}, 10)
```

## 使用示例

### 完整示例代码
```go
func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 初始化服务
    service, err := NewMeiliSearchService(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 准备数据
    books := []Book{
        {
            ID:      "1",
            Title:   "Go语言高级编程",
            Authors: []string{"张三", "李四"},
            Price:   59.9,
            Genres:  []string{"编程", "计算机"},
            Publish: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
        },
        // 更多书籍...
    }

    // 添加文档
    taskID, err := service.AddDocuments("books", books)
    if err != nil {
        log.Fatal(err)
    }

    // 等待任务完成
    if err := service.WaitForTask(taskID, 10*time.Second); err != nil {
        log.Fatal(err)
    }

    // 执行搜索
    results, err := service.Search("books", "编程", "price < 70", []string{"price:asc"}, 10)
    if err != nil {
        log.Fatal(err)
    }

    // 处理结果
    for _, hit := range results {
        fmt.Printf("%+v\n", hit)
    }
}
```
