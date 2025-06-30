# Supabase Auth - 阿里云短信 & 华为云短信增强版

本项目基于 [Supabase Auth](https://github.com/supabase/auth) 进行扩展，添加了**阿里云短信服务** 和 **华为云短信服务**的功能。

## 🚀 新增功能

### 1. 阿里云短信服务 (Aliyun SMS)
- ✅ 完整的阿里云短信 API 集成
- ✅ 支持 HMAC-SHA1 签名验证
- ✅ 支持中文短信签名
- ✅ 支持 OTP 验证码发送
- ✅ 完整的错误处理和响应解析

### 2. 华为云短信服务 (HuaweiCloud SMS)
- ✅ 添加华为云短信 API 集成
- ✅ 完整的 `VerifyOTP` 方法实现
- ✅ 完善的接口功能支持

## 📁 文件修改清单

### 新增文件
- `internal/api/sms_provider/aliyun.go` - 阿里云短信提供者实现

### 修改文件
- `internal/conf/configuration.go` - 添加阿里云短信配置结构体
- `internal/api/sms_provider/sms_provider.go` - 注册阿里云和华为云短信提供者
- `internal/api/sms_provider/huaweicloud.go` - 添加华为云SMS完整功能支持

## 🔧 配置说明

### 环境变量配置

#### 基础SMS配置
```bash
# 短信提供商选择
SMS_PROVIDER=aliyun

# 短信模板代码
SMS_TEMPLATE=SMS_154950909
```

#### 阿里云短信配置
```bash
# 阿里云 Access Key ID
SMS_ALIYUN_ACCESS_KEY_ID=your_access_key_id

# 阿里云 Access Key Secret  
SMS_ALIYUN_ACCESS_KEY_SECRET=your_access_key_secret

# 阿里云短信服务终端
SMS_ALIYUN_ENDPOINT=https://dysmsapi.aliyuncs.com

# 阿里云短信签名（支持中文）
SMS_ALIYUN_SIGN_NAME=您的短信签名

# 阿里云短信扩展码（可选）
SMS_ALIYUN_SMS_UP_EXTEND_CODE=
```

#### GoTrue 环境变量映射
```bash
# 在 docker-compose.yml 中映射到 GoTrue 容器
GOTRUE_SMS_ALIYUN_ACCESS_KEY_ID=${SMS_ALIYUN_ACCESS_KEY_ID}
GOTRUE_SMS_ALIYUN_ACCESS_KEY_SECRET=${SMS_ALIYUN_ACCESS_KEY_SECRET}
GOTRUE_SMS_ALIYUN_ENDPOINT=${SMS_ALIYUN_ENDPOINT}
GOTRUE_SMS_ALIYUN_SIGN_NAME=${SMS_ALIYUN_SIGN_NAME}
GOTRUE_SMS_ALIYUN_SMS_UP_EXTEND_CODE=${SMS_ALIYUN_SMS_UP_EXTEND_CODE}
```

## 🏗️ Docker 构建

### 构建自定义镜像
```bash
# 在 auth 目录下执行
cd auth
sudo docker build -t supabase/gotrue:local .
```

### 在 Supabase 项目中使用
确保 `docker-compose.yml` 中的 auth 服务使用自定义镜像：
```yaml
services:
  auth:
    image: supabase/gotrue:local
    # ... 其他配置
```

## 🧪 测试使用

### 1. 短信注册测试
```python
# 使用提供的测试脚本
./run_sms_test.sh register +8613800138000

# 或直接使用 Python
python test_sms_auth.py register +8613800138000
```

### 2. OTP 登录测试  
```python
# 输入收到的验证码
./run_sms_test.sh login +8613800138000 123456

# 或直接使用 Python
python test_sms_auth.py login +8613800138000 123456
```

## 📊 API 集成详解

### 阿里云短信 API 实现

#### 关键特性
1. **完整签名验证**：实现 HMAC-SHA1 签名算法
2. **参数编码**：正确处理 URL 编码和参数排序
3. **中文支持**：支持中文短信签名和模板内容
4. **错误处理**：完整的 API 错误响应解析

#### 签名算法
```go
// 1. 参数排序
// 2. 构建查询字符串
// 3. 构建待签名字符串：POST&%2F&{encodedQuery}
// 4. HMAC-SHA1 签名并 Base64 编码
signature := generateSignature(params, accessKeySecret)
```

#### API 调用流程
```go
// 1. 构建请求参数
params := map[string]string{
    "Action": "SendSms",
    "Version": "2017-05-25",
    "AccessKeyId": accessKeyId,
    "PhoneNumbers": phone,
    "SignName": signName,
    "TemplateCode": templateCode,
    "TemplateParam": `{"code":"123456"}`,
    // ... 其他参数
}

// 2. 生成签名
signature := generateSignature(params, accessKeySecret)

// 3. 发送 POST 请求到阿里云 API
// 4. 解析 JSON 响应
```

## 🔍 故障排除

### 常见问题

#### 1. 签名验证失败
```
错误：SignatureDoesNotMatch
解决：检查 AccessKeySecret 是否正确，确保时间戳格式为 UTC
```

#### 2. 模板参数错误
```
错误：InvalidTemplateCode
解决：确保模板代码格式正确，如 SMS_154950909
```

#### 3. 中文签名问题
```
错误：InvalidSignName  
解决：确保签名名称与阿里云控制台完全一致，无需加引号
```

#### 4. 权限错误
```
错误：Forbidden.RAM
解决：确保 RAM 用户有短信发送权限
```

### 调试方法

#### 查看 Auth 服务日志
```bash
# 查看 Docker 容器日志
sudo docker logs supabase-auth | tail -20

# 实时监控日志
sudo docker logs -f supabase-auth
```

#### 检查环境变量
```bash
# 在容器内检查环境变量
sudo docker exec supabase-auth env | grep ALIYUN
```

## 📈 性能优化

### 1. 连接池优化
- 使用 HTTP 客户端连接池
- 设置合理的超时时间

### 2. 错误重试
- 实现指数退避重试机制
- 区分可重试和不可重试错误

### 3. 监控指标
- 短信发送成功率
- API 响应时间
- 错误码分布

## 🛡️ 安全考虑

### 1. 密钥管理
- 使用环境变量存储敏感信息
- 定期轮换 AccessKey
- 最小权限原则

### 2. 防刷机制
- 实现频率限制
- IP 白名单
- 验证码有效期

### 3. 日志脱敏
- 避免记录敏感参数
- 脱敏手机号显示

## 🔄 升级指南

### 从原版 Supabase Auth 升级

1. **备份现有配置**
2. **添加阿里云配置项**到环境变量
3. **重新构建 Docker 镜像**
4. **更新 docker-compose.yml**
5. **测试短信功能**

### 版本兼容性
- ✅ 与原版 Supabase Auth API 完全兼容
- ✅ 支持现有的邮件和其他 SMS 提供商
- ✅ 向后兼容所有现有功能

## 📞 技术支持

### 联系方式
- GitHub Issues: [提交问题](https://github.com/estel-li/supabase_auth_aliyun_wechat/issues)
- 邮箱：请通过 GitHub 联系

### 贡献指南
欢迎提交 Pull Request 来改进功能：
1. Fork 本仓库
2. 创建功能分支
3. 提交更改
4. 发起 Pull Request

## 📜 许可证

本项目遵循 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

---

**感谢使用！如有问题请提交 Issue。** 🙏 