# OCR 识别服务文档

## 目录
1. [概述](#概述)
2. [核心结构](#核心结构)
   - [TradeNoResult](#tradenoresult)
   - [TradeNoParserService](#tradenoparserservice)
3. [接口定义](#接口定义)
   - [TradeNoParser](#tradenoparser)
4. [具体实现](#具体实现)
   - [TaobaoTraderParser](#taobaotraderparser)
5. [使用示例](#使用示例)
6. [错误处理](#错误处理)
6. [完整代码](#完整代码)

## 概述

本服务使用了go版本的ocr服务gosseract，提供从图片中提取淘宝交易信息的功能，包括：
- 订单编号
- 订单状态
- 交易金额
- 原始OCR文本

访问[gosseract](https://github.com/otiai10/gosseract?tab=readme-ov-file)了解更多安装方法，以及相关的API文档

## 核心结构

### TradeNoResult

```go
type TradeNoResult struct {
    TradeNo     string  `json:"trade_no"`           // 交易流水号
    OrderStatus string  `json:"order_status"`       // 订单状态
    Amount      float64 `json:"amount"`             // 金额
    RawText     string  `json:"raw_text,omitempty"` // 原始OCR文本
}
```

### TradeNoParserService

```go
type TradeNoParserService struct {
    parser TradeNoParser // 依赖的解析器（接口类型）
}
```

## 接口定义

### TradeNoParser

```go
type TradeNoParser interface {
    ParseImage(filePath string) (TradeNoResult, error)
}
```

## 具体实现

### TaobaoTraderParser

淘宝交易信息解析器的具体实现：

```go
type TaobaoTraderParser struct {
    Request *request.MultipartRequest
    Logger  *logx.Logger
}
```

#### 主要方法

##### ParseImage

```go
func (s *TaobaoTraderParser) ParseImage(filePath string) (result TradeNoResult, err error)
```

功能：解析指定图片文件，提取交易信息

参数：
- `filePath`: 要解析的图片文件路径

返回值：
- `result`: 解析结果(TradeNoResult)
- `err`: 错误信息

处理流程：
1. 验证文件路径
2. 调用OCR服务上传图片
3. 解析返回结果
4. 提取订单编号、状态和金额信息

##### parseAmount

```go
func (s *TaobaoTraderParser) parseAmount(text string) float64
```

内部方法，用于解析金额字符串

## 使用示例

```go
// 创建请求对象和日志对象
req := request.NewMultipartRequest()
logger := logx.NewLogger()

// 创建解析器实例
parser := NewTaobaoTraderParser(req, logger)

// 创建服务并设置解析器
service := NewTradeNoParserService().SetTradeOperationObject(parser)

// 调用解析服务
result, err := service.Parse("/path/to/image.jpg")
if err != nil {
    // 处理错误
    fmt.Printf("解析失败: %v\n", err)
    return
}

// 使用解析结果
fmt.Printf("订单编号: %s\n", result.TradeNo)
fmt.Printf("订单状态: %s\n", result.OrderStatus)
fmt.Printf("交易金额: %.2f\n", result.Amount)
```
## 完整代码

### trade_no.go
```go
package ocr

import (
	"errors"
	"fmt"
)

type (
	// 识别结果
	TradeNoResult struct {
		TradeNo     string  `json:"trade_no"`           // 交易流水号
		OrderStatus string  `json:"order_status"`       // 订单状态
		Amount      float64 `json:"amount"`             // 金额
		RawText     string  `json:"raw_text,omitempty"` // 原始OCR文本
	}

	// 识别服务
	TradeNoParserService struct {
		parser TradeNoParser // 依赖的解析器（接口类型）
	}
)

// 解析器接口（命名更清晰）
type TradeNoParser interface {
	ParseImage(filePath string) (TradeNoResult, error)
}

// 构造函数：初始化空服务（后续需调用 SetTradeOperationObject）
func NewTradeNoParserService() *TradeNoParserService {
	return &TradeNoParserService{} // parser 默认为 nil
}

// 动态设置解析器（支持链式调用）
func (s *TradeNoParserService) SetTradeOperationObject(parser TradeNoParser) *TradeNoParserService {
	s.parser = parser
	return s
}

// 识别逻辑
func (s *TradeNoParserService) Parse(filePath string) (TradeNoResult, error) {
	if s.parser == nil {
		return TradeNoResult{}, errors.New("未设置识别解析器，请调用 SetTradeOperationObject()")
	}
	result, err := s.parser.ParseImage(filePath)
	if err != nil {
		return TradeNoResult{}, fmt.Errorf("识别失败: %w", err) // 包装错误
	}
	return result, nil
}

```

### taobao_trader_parser.go
```go
package ocr

import (
	"cms/service/request"
	"cms/util/errorx"
	"cms/util/logx"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type TaobaoTraderParser struct {
	Request *request.MultipartRequest
	Logger  *logx.Logger
}

func NewTaobaoTraderParser(r *request.MultipartRequest, logger *logx.Logger) *TaobaoTraderParser {
	return &TaobaoTraderParser{
		Request: r,
		Logger:  logger,
	}
}

// 解析图片
//
// filePath 图片路径
//
// 返回值：
//
// body 原文解析结果   result 订单信息  err 错误信息
func (s *TaobaoTraderParser) ParseImage(filePath string) (result TradeNoResult, err error) {
	if filePath == "" {
		return TradeNoResult{}, errorx.New("文件路径不能为空")
	}

	// 上传图片
	responseData, err := s.Request.AddFile("file", filePath).Send(http.MethodPost, "http://127.0.0.1:8080/ocr/file?languages=chi_sim")

	if err != nil {
		return
	}

	s.Logger.Log.Info("解析图片返回结果：", responseData.Result)

	if responseData.Result == "" {
		return result, errorx.Wrap(err, fmt.Sprintf("图片OCR识别失败，文件路径: %s", filePath))
	}

	result.RawText = responseData.Result

	// 正则表达式匹配订单编号
	orderIDRegex := regexp.MustCompile(`订单编号\s+(\d+)`)
	orderIDMatch := orderIDRegex.FindStringSubmatch(responseData.Result)
	if len(orderIDMatch) > 1 {
		result.TradeNo = orderIDMatch[1]
	}

	// 正则表达式匹配订单状态
	orderStatusRegex := regexp.MustCompile(`(已揽件|运输中|派送中|待取件|已签收|待付款|已发货|已付款|已取消|已退款|交易成功|退款中|交易关闭)`)
	orderStatusMatch := orderStatusRegex.FindStringSubmatch(responseData.Result)
	if len(orderStatusMatch) > 1 {
		result.OrderStatus = orderStatusMatch[1]
	}

	// 正则表达式匹配实付款金额
	paymentRegex := regexp.MustCompile(`实付款\s+Y#(\d+(?:\.\d+)?)`)
	paymentMatch := paymentRegex.FindStringSubmatch(responseData.Result)
	if len(paymentMatch) > 1 {
		result.Amount += s.parseAmount(paymentMatch[1])
	}

	// 正则表达式匹配确认收货后再付款金额
	postPaymentRegex := regexp.MustCompile(`确认收货后再付款\s+(\d+(?:\.\d+)?)`)
	postPaymentMatch := postPaymentRegex.FindStringSubmatch(responseData.Result)
	if len(postPaymentMatch) > 1 {
		result.Amount += s.parseAmount(postPaymentMatch[1])
	}

	return result, nil
}

// 金额解析函数
func (s *TaobaoTraderParser) parseAmount(text string) float64 {
	amount, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0 // 或记录错误日志
	}
	return amount
}

```