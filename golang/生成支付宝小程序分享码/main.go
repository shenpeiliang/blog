package main

import (
	"demo/alipay"
	"fmt"
)

func main() {
	service := alipay.NewAlipayOpenAppQrcodeService(
		"xxxxxxxxxxxxxx",
		"./keys/app_private_key.pem",
		"./keys/alipay_public_key.pem",
	)
	result, err := service.Create(alipay.AlipayOpenAppQrcodeCreateRequest{
		UrlParam:   "page/component/component-pages/view/view", // 小程序的页面路径
		QueryParam: "key1=value1&key2=value2",                  // 小程序的启动参数
		Describe:   "二维码描述",                                    // 码描述
		Color:      "0x00BFFF",                                 // 码颜色（可选）
		Size:       "s",                                        // 合成后图片的大小（可选）
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
