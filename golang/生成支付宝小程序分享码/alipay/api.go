package alipay

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	//支付宝网关地址
	GatewayURL = "https://openapi.alipay.com"

	//创建小程序二维码
	AlipayOpenAppQrcodeCreateURL = "/v3/alipay/open/app/qrcode/create"
)

type (
	AlipayErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	// 请求参数结构
	AlipayOpenAppQrcodeCreateRequest struct {
		UrlParam   string `json:"url_param"`       // 小程序的页面路径
		QueryParam string `json:"query_param"`     // 小程序的启动参数
		Describe   string `json:"describe"`        // 码描述
		Color      string `json:"color,omitempty"` // 码颜色（可选）
		Size       string `json:"size,omitempty"`  // 合成后图片的大小（可选）
	}

	// 响应参数结构
	AlipayOpenAppQrcodeCreateResponse struct {
		AlipayErrorResponse
		QrCodeUrl            string `json:"qr_code_url"`              // 方形二维码图片链接地址
		QrCodeUrlCircleWhite string `json:"qr_code_url_circle_white"` // 白底圆码链接地址
		QrCodeUrlCircleBlue  string `json:"qr_code_url_circle_blue"`  // 蓝底圆码链接地址
	}

	AlipayOpenAppQrcodeService struct {
		AppID               string
		AppPrivateKeyPath   string
		AlipayPublicKeyPath string
	}
)

// 创建小程序二维码服务
func NewAlipayOpenAppQrcodeService(appID, appPrivateKeyPath, alipayPublicKeyPath string) *AlipayOpenAppQrcodeService {
	return &AlipayOpenAppQrcodeService{
		AppID:               appID,
		AppPrivateKeyPath:   appPrivateKeyPath,
		AlipayPublicKeyPath: alipayPublicKeyPath,
	}
}

// 创建小程序二维码
func (s *AlipayOpenAppQrcodeService) Create(query AlipayOpenAppQrcodeCreateRequest) (result *AlipayOpenAppQrcodeCreateResponse, err error) {
	//请求参数序列化
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	//请求参数签名
	authorization, nonce, err := getAuthorizationContent(SignContentParam{
		AppID:             s.AppID,
		Method:            RequestMethodPost,
		URI:               AlipayOpenAppQrcodeCreateURL,
		BodyJson:          string(body),
		AppPrivateKeyPath: s.AppPrivateKeyPath,
	})

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":      "application/json",
		"alipay-request-id": nonce,
		"authorization":     authorization,
	}

	//发送请求
	response, err := NewRequester(
		WithAlipayPublicKeyPath(s.AlipayPublicKeyPath),
	).Post(GatewayURL+AlipayOpenAppQrcodeCreateURL, headers, body)

	if err != nil {
		return nil, fmt.Errorf("request alipay open app qrcode create failed, err: %v", err)
	}
	//响应参数反序列化
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}

	if result.Code != "" {
		return nil, errors.New(result.Message)
	}
	return result, nil
}
