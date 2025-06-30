package sms_provider

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/supabase/auth/internal/conf"
)

type AliyunProvider struct {
	Config *conf.AliyunProviderConfiguration
}

type AliyunSmsResponse struct {
	BizId     string `json:"BizId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
}

// SendMessage implements SmsProvider.
func (p *AliyunProvider) SendMessage(phone string, message string, channel string, otp string) (string, error) {
	switch channel {
	case SMSProvider:
		return p.SendSms(phone, message, otp)
	default:
		return "", fmt.Errorf("channel type %q is not supported for Aliyun SMS service", channel)
	}
}

func (p *AliyunProvider) SendSms(phone, templateCode, otp string) (string, error) {
	// 构建模板参数，将 OTP 作为模板变量
	templateParam := fmt.Sprintf(`{"code":"%s"}`, otp)

	// 构建阿里云API请求参数
	params := map[string]string{
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"AccessKeyId":      p.Config.AccessKeyId,
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   strconv.FormatInt(time.Now().UnixNano(), 10),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"PhoneNumbers":     phone,
		"SignName":         p.Config.SignName,
		"TemplateCode":     templateCode,
		"TemplateParam":    templateParam,
	}

	// 添加可选的扩展码
	if p.Config.SmsUpExtendCode != "" {
		params["SmsUpExtendCode"] = p.Config.SmsUpExtendCode
	}

	// 生成签名
	signature := p.generateSignature(params, p.Config.AccessKeySecret)
	params["Signature"] = signature

	// 构建POST请求
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	// 发送请求
	req, err := http.NewRequest("POST", p.Config.Endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send SMS via Aliyun: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var smsResp AliyunSmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 检查响应状态
	if smsResp.Code != "OK" {
		return "", fmt.Errorf("Aliyun SMS error: Code=%s, Message=%s", smsResp.Code, smsResp.Message)
	}

	return smsResp.BizId, nil
}

// generateSignature 生成阿里云API签名
func (p *AliyunProvider) generateSignature(params map[string]string, accessKeySecret string) string {
	// 1. 对参数进行排序
	var keys []string
	for k := range params {
		if k != "Signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 2. 构建查询字符串
	var query []string
	for _, k := range keys {
		query = append(query, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	queryString := strings.Join(query, "&")

	// 3. 构建待签名字符串
	stringToSign := "POST&%2F&" + url.QueryEscape(queryString)

	// 4. 计算签名
	key := accessKeySecret + "&"
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}

// VerifyOTP implements SmsProvider.
func (p *AliyunProvider) VerifyOTP(phone string, otp string) error {
	// 阿里云短信服务不支持直接的 OTP 验证
	// 这个方法是接口要求的，但实际上不实现
	return fmt.Errorf("OTP verification not supported by Aliyun SMS provider")
}

func NewAliyunProvider(config conf.AliyunProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &AliyunProvider{
		Config: &config,
	}, nil
}
