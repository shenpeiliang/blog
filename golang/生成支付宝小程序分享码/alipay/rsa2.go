package alipay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// 读取指定文件的内容
func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// 读取私钥
func ReadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	//读取私钥文件内容
	privateKey, err := ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	//解析私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("系统配置出错")
	}

	privateKeyInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyInterface, nil
}

// 读取公钥
func ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	//读取公钥文件内容
	publicKey, err := ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	//解析公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, fmt.Errorf("系统配置出错")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKeyInterface.(*rsa.PublicKey), nil
}
