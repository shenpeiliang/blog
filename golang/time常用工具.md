# TimeFormatter 时间格式化工具文档

## 概述

`TimeFormatter` 是一个 Go 语言的时间处理工具包，提供了一系列方便的时间格式化和计算功能，支持时区设置。

## 功能特性

- **时区支持**：所有时间操作都支持指定时区
- **时间格式化**：支持自定义格式输出时间
- **时间解析**：将字符串解析为时间戳
- **时间计算**：提供日期加减、当天起止时间等功能
- **人性化时间描述**：生成"刚刚"、"X分钟前"等友好时间格式

## 安装

```go
import "your-module-path/timex"
```

## 核心结构体

### TimeFormatter

```go
type TimeFormatter struct {
    Location *time.Location
}
```

- `Location`：时区设置，默认为初始化时传入的时区

## 初始化

### NewTimeFormatter

```go
func NewTimeFormatter(loc *time.Location) *TimeFormatter
```

示例：
```go
// 使用本地时区
tf := NewTimeFormatter(time.Local)

// 使用UTC时区
tf := NewTimeFormatter(time.UTC)

// 使用上海时区
loc, _ := time.LoadLocation("Asia/Shanghai")
tf := NewTimeFormatter(loc)
```

## 方法说明

### 1. 时间格式化

#### Format

```go
func (tf *TimeFormatter) Format(sec uint, layout ...string) string
```

- `sec`：Unix 时间戳（秒）
- `layout`：可选，时间格式字符串，默认为 "2006-01-02 15:04:05"

示例：
```go
now := tf.UnixTime()
fmt.Println(tf.Format(now)) // 2023-07-20 14:30:00
fmt.Println(tf.Format(now, "2006/01/02")) // 2023/07/20
```

### 2. 时间解析

#### ParseTime

```go
func (tf *TimeFormatter) ParseTime(value string, layout ...string) (ut uint, err error)
```

- `value`：时间字符串
- `layout`：可选，时间格式字符串，默认为 "2006-01-02 15:04:05"

示例：
```go
ts, err := tf.ParseTime("2023-07-20 14:30:00")
if err != nil {
    // 处理错误
}
```

### 3. 获取当前时间

#### Now

```go
func (tf *TimeFormatter) Now() time.Time
```

返回当前时间（带时区）

#### UnixTime

```go
func (tf *TimeFormatter) UnixTime() uint
```

返回当前 Unix 时间戳（秒）

### 4. 时间计算

#### BeforeDay / AfterDay

```go
func (tf *TimeFormatter) BeforeDay(day int) uint
func (tf *TimeFormatter) AfterDay(day int) uint
```

计算指定天数前/后的时间戳

示例：
```go
yesterday := tf.BeforeDay(1) // 昨天此时的时间戳
tomorrow := tf.AfterDay(1)   // 明天此时的时间戳
```

#### StartOfDay / EndOfDay

```go
func (tf *TimeFormatter) StartOfDay() uint
func (tf *TimeFormatter) EndOfDay() uint
```

获取当天开始（00:00:00）和结束（23:59:59）的时间戳

### 5. 人性化时间描述

#### TimeDescription

```go
func (tf *TimeFormatter) TimeDescription(t uint) string
```

生成易读的时间描述，如：
- "刚刚"
- "5分钟前"
- "3小时前"
- "昨天 14:30"
- "前天 10:15"
- "3天前 09:00"
- "07-15 14:30"（同年）
- "2022-12-01 08:00"（跨年）

## 使用示例

### 基本使用

```go
package main

import (
	"fmt"
	"time"
	"timex"
)

func main() {
	// 初始化（使用上海时区）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tf := timex.NewTimeFormatter(loc)

	// 当前时间
	now := tf.UnixTime()
	fmt.Println("当前时间:", tf.Format(now))

	// 昨天此时
	yesterday := tf.BeforeDay(1)
	fmt.Println("昨天此时:", tf.Format(yesterday))

	// 解析时间
	if ts, err := tf.ParseTime("2023-07-20 14:30:00"); err == nil {
		fmt.Println("解析后的时间戳:", ts)
	}

	// 人性化时间描述
	fmt.Println(tf.TimeDescription(now - 30))    // 刚刚
	fmt.Println(tf.TimeDescription(now - 300))   // 5分钟前
	fmt.Println(tf.TimeDescription(now - 7200))  // 2小时前
	fmt.Println(tf.TimeDescription(tf.StartOfDay() - 1)) // 昨天 23:59
}
```

## 完整代码

```go
package timex

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type TimeFormatter struct {
	Location *time.Location
}

// 新进一个时间格式化器
//
// loc时区 例如：time.Local（默认）、time.UTC、 time.LoadLocation("Asia/Shanghai")
func NewTimeFormatter(loc *time.Location) *TimeFormatter {
	return &TimeFormatter{Location: loc}
}

// 格式化时间
func (tf *TimeFormatter) Format(sec uint, layout ...string) string {
	l := "2006-01-02 15:04:05"

	if len(layout) > 0 {
		l = layout[0]
	}

	return time.Unix(int64(sec), 0).In(tf.Location).Format(l)
}

// 日期转时间戳
func (tf *TimeFormatter) ParseTime(value string, layout ...string) (ut uint, err error) {
	l := "2006-01-02 15:04:05"

	if len(layout) > 0 {
		l = layout[0]
	}

	t, err := time.ParseInLocation(l, value, tf.Location)
	if err != nil {
		err = errors.New("日期解析错误")
		return
	}

	ut = uint(t.Unix())

	return
}

// 当前时间
func (tf *TimeFormatter) Now() time.Time {
	return time.Now().In(tf.Location)
}

// 当前时间戳
func (tf *TimeFormatter) UnixTime() uint {
	return uint(tf.Now().Unix())
}

// 指定天数前的当前时间戳
func (tf *TimeFormatter) BeforeDay(day int) uint {
	t := tf.Now().AddDate(0, 0, -day) // 使用 AddDate 安全计算时间
	return uint(t.Unix())
}

// 指定天数后的当前时间戳
func (tf *TimeFormatter) AfterDay(day int) uint {
	t := tf.Now().AddDate(0, 0, day) // 使用 AddDate 安全计算时间
	return uint(t.Unix())
}

// 当天凌晨时间戳
func (tf *TimeFormatter) StartOfDay() uint {
	t := tf.Now()
	return uint(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix())
}

// 当天结束时间戳
func (tf *TimeFormatter) EndOfDay() uint {
	t := tf.Now()
	return uint(time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix())
}

// 格式化时间描述
//
// t时间戳
func (tf *TimeFormatter) TimeDescription(t uint) string {
	//当前时间
	now := tf.Now()

	//时间戳转时间
	tTime := time.Unix(int64(t), 0).In(tf.Location)

	//当前时间和t的差值
	timeRange := now.Unix() - tTime.Unix()

	if timeRange < 60 {
		return "刚刚"
	} else if timeRange < 3600 {
		return fmt.Sprintf("%d分钟前", uint(math.Ceil(float64(timeRange)/60)))
	} else if timeRange < 3600*12 {
		return fmt.Sprintf("%d小时前", uint(math.Ceil(float64(timeRange)/3600)))
	}

	// 计算天数差
	dayRange := int(now.Sub(tTime).Hours() / 24)
	if now.YearDay() == tTime.YearDay() && now.Year() == tTime.Year() {
		return fmt.Sprintf("今天 %s", tTime.Format("15:04"))
	} else if now.YearDay()-1 == tTime.YearDay() && now.Year() == tTime.Year() {
		return fmt.Sprintf("昨天 %s", tTime.Format("15:04"))
	} else if now.YearDay()-2 == tTime.YearDay() && now.Year() == tTime.Year() {
		return fmt.Sprintf("前天 %s", tTime.Format("15:04"))
	} else if dayRange > 2 && dayRange < 16 {
		return fmt.Sprintf("%d天前 %s", dayRange, tTime.Format("15:04"))
	} else if now.Year() == tTime.Year() {
		return tTime.Format("01-02 15:04")
	}

	return tTime.Format("2006-01-02 15:04")

}

```