# 使用Go生成支付宝小程序分享码指南

## 概述

本文档介绍如何使用Go语言SDK生成支付宝小程序分享二维码。支付宝小程序二维码可以帮助用户快速访问小程序特定页面，并携带自定义参数。

## 准备工作

在开始之前，请确保您已经：

1. 拥有支付宝开放平台开发者账号
2. 创建了小程序应用并获取了AppID
3. 准备好了RSA密钥对（应用私钥和支付宝公钥）

## 初始化服务

首先需要初始化二维码生成服务：

```go
service := alipay.NewAlipayOpenAppQrcodeService(
    "xxxxxxxxxxxxxx",              // 您的支付宝小程序AppID
    "./keys/app_private_key.pem",    // 应用私钥文件路径
    "./keys/alipay_public_key.pem",  // 支付宝公钥文件路径
)
```

## 生成二维码

使用`Create`方法生成小程序二维码：

```go
result, err := service.Create(alipay.AlipayOpenAppQrcodeCreateRequest{
    UrlParam:   "page/component/component-pages/view/view", // 小程序的页面路径
    QueryParam: "key1=value1&key2=value2",                  // 小程序的启动参数
    Describe:   "二维码描述",                                   // 码描述
    Color:      "0x00BFFF",                                 // 码颜色（可选）
    Size:       "s",                                        // 合成后图片的大小（可选）
})

if err != nil {
    fmt.Println("生成二维码失败:", err)
    return
}

fmt.Println("二维码生成结果:", result)
```

## 参数说明

| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| UrlParam | string | 是 | 小程序中能访问到的页面路径 |
| QueryParam | string | 否 | 小程序的启动参数，格式为key1=value1&key2=value2 |
| Describe | string | 否 | 对二维码的描述 |
| Color | string | 否 | 二维码颜色，格式为十六进制颜色值，如0x00BFFF |
| Size | string | 否 | 合成后图片的大小，可选值：s(小)、m(中)、l(大) |

## 签名验证

签名过程使用SHA256withRSA算法：

```go
// 获取签名字符串
func getAuthString(param AuthorizationParam) string {
	if param.AppCertSn == "" {
		return fmt.Sprintf("app_id=%s,nonce=%s,timestamp=%s", param.AppID, param.Nonce, param.Timestamp)
	}
	return fmt.Sprintf("app_id=%s,app_cert_sn=%s,nonce=%s,timestamp=%s", param.AppID, param.AppCertSn, param.Nonce, param.Timestamp)
}

// 获取签名内容
func getAuthorizationContent(param SignContentParam) (string, string, error) {
	// 解析私钥
	privateKey, err := ReadPrivateKey(param.AppPrivateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("Failed to read private key: %v", err)
	}

	// 生成签名
	timestamp := time.Now().Format("20060102150405")
	nonce := strings.ReplaceAll(fmt.Sprintf("%d", time.Now().UnixNano()), "-", "")

	// 签名字符串
	authString := getAuthString(AuthorizationParam{
		AppID:     param.AppID,
		Nonce:     nonce,
		Timestamp: timestamp,
	})

	// 签名
	signString, err := Sign(Signature{
		AuthString: authString,
		Method:     param.Method,
		URI:        param.URI,
		BodyJson:   param.BodyJson,
	}, privateKey)

	if err != nil {
		return "", "", fmt.Errorf("Failed to sign request: %v", err)
	}

	return SignTypeRSA + " " + authString + ",sign=" + signString, nonce, nil
}

// 对数据进行签名
func Sign(param Signature, privateKey *rsa.PrivateKey) (string, error) {
	// 构造待签名内容
	content := param.AuthString + "\n" + param.Method + "\n" + param.URI + "\n" + param.BodyJson + "\n"

	// 计算哈希值
	hash := sha256.New()
	hash.Write([]byte(content))
	hashed := hash.Sum(nil)

	// 生成签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

```

验证签名
```go
func verifySignByCert(res *http.Response, body []byte, alipayPublicKeyPath string) (err error) {
	ts := res.Header.Get(HeaderTimestamp)
	nonce := res.Header.Get(HeaderNonce)
	sign := res.Header.Get(HeaderSignature)

	signData := ts + "\n" + nonce + "\n" + string(body) + "\n"

	// 解析公钥
	publicKey, err := ReadPublicKey(alipayPublicKeyPath)
	if err != nil {
		return fmt.Errorf("Failed to read public key: %v", err)
	}

	// 验证签名
	return Verify([]byte(signData), sign, publicKey)
}
```

## 示例代码

完整示例：

```go
package main

import (
	"fmt"
	"github.com/yourpackage/alipay"
)

func main() {
	// 初始化服务
	service := alipay.NewAlipayOpenAppQrcodeService(
		"2021004145619818",
		"./keys/app_private_key.pem",
		"./keys/alipay_public_key.pem",
	)

	// 生成二维码
	result, err := service.Create(alipay.AlipayOpenAppQrcodeCreateRequest{
		UrlParam:   "page/component/component-pages/view/view",
		QueryParam: "key1=value1&key2=value2",
		Describe:   "商品详情页二维码",
		Color:      "0x00BFFF",
		Size:       "m",
	})

	if err != nil {
		fmt.Println("生成二维码失败:", err)
		return
	}

	fmt.Println("二维码URL:", result.QrCodeUrl)
	fmt.Println("二维码图片:", result.QrCodeImage)
}
```