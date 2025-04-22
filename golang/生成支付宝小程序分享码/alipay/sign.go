package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	//签名方法
	SignTypeRSA = "ALIPAY-SHA256withRSA"

	//请求头字段名
	HeaderTimestamp = "Alipay-Timestamp"
	HeaderNonce     = "Alipay-Nonce"
	HeaderSignature = "Alipay-Signature"
)

type (
	Signature struct {
		AuthString string
		Method     string
		URI        string
		BodyJson   string
	}

	AuthorizationParam struct {
		AppID     string
		AppCertSn string // 应用公钥证书序列号
		Nonce     string
		Timestamp string
	}

	SignContentParam struct {
		AppID             string
		Method            string
		URI               string
		BodyJson          string
		AppPrivateKeyPath string
	}
)

// 自动验证签名
//
// 参考文档：https://opendocs.alipay.com/open-v3/054d0z?pathHash=dcad8d5c
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

// 验证签名
func Verify(data []byte, sign string, publicKey *rsa.PublicKey) error {
	// 解码base64字符串
	decodedSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(data)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodedSign)
	if err != nil { //签名验证失败
		return err
	}

	return nil
}
