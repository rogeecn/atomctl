package errorx

import "net/http"

// 预定义错误 - 数据相关
var (
    ErrRecordNotFound     = NewError(CodeRecordNotFound, http.StatusNotFound, "记录不存在")
    ErrRecordDuplicated   = NewError(CodeRecordDuplicated, http.StatusConflict, "记录重复")
    ErrDataCorrupted      = NewError(CodeDataCorrupted, http.StatusBadRequest, "数据损坏")
    ErrDataTooLarge       = NewError(CodeDataTooLarge, http.StatusRequestEntityTooLarge, "数据过大")
    ErrDataValidationFail = NewError(CodeDataValidationFail, http.StatusBadRequest, "数据验证失败")
    ErrConstraintViolated = NewError(CodeConstraintViolated, http.StatusConflict, "约束违规")
    ErrDataExpired        = NewError(CodeDataExpired, http.StatusGone, "数据已过期")
    ErrDataLocked         = NewError(CodeDataLocked, http.StatusLocked, "数据已锁定")
)

// 预定义错误 - 请求相关
var (
    ErrBadRequest        = NewError(CodeBadRequest, http.StatusBadRequest, "请求错误")
    ErrMissingParameter  = NewError(CodeMissingParameter, http.StatusBadRequest, "缺少必需参数")
    ErrInvalidParameter  = NewError(CodeInvalidParameter, http.StatusBadRequest, "参数无效")
    ErrParameterTooLong  = NewError(CodeParameterTooLong, http.StatusBadRequest, "参数过长")
    ErrParameterTooShort = NewError(CodeParameterTooShort, http.StatusBadRequest, "参数过短")
    ErrInvalidFormat     = NewError(CodeInvalidFormat, http.StatusBadRequest, "格式无效")
    ErrUnsupportedMethod = NewError(CodeUnsupportedMethod, http.StatusMethodNotAllowed, "不支持的请求方法")
    ErrRequestTooLarge   = NewError(CodeRequestTooLarge, http.StatusRequestEntityTooLarge, "请求体过大")
    ErrInvalidJSON       = NewError(CodeInvalidJSON, http.StatusBadRequest, "JSON格式错误")
    ErrInvalidXML        = NewError(CodeInvalidXML, http.StatusBadRequest, "XML格式错误")
)

// 预定义错误 - 认证授权
var (
    ErrUnauthorized       = NewError(CodeUnauthorized, http.StatusUnauthorized, "未授权")
    ErrForbidden          = NewError(CodeForbidden, http.StatusForbidden, "禁止访问")
    ErrTokenExpired       = NewError(CodeTokenExpired, http.StatusUnauthorized, "Token已过期")
    ErrTokenInvalid       = NewError(CodeTokenInvalid, http.StatusUnauthorized, "Token无效")
    ErrTokenMissing       = NewError(CodeTokenMissing, http.StatusUnauthorized, "Token缺失")
    ErrPermissionDenied   = NewError(CodePermissionDenied, http.StatusForbidden, "权限不足")
    ErrAccountDisabled    = NewError(CodeAccountDisabled, http.StatusForbidden, "账户已禁用")
    ErrAccountLocked      = NewError(CodeAccountLocked, http.StatusLocked, "账户已锁定")
    ErrInvalidCredentials = NewError(CodeInvalidCredentials, http.StatusUnauthorized, "凭据无效")
    ErrSessionExpired     = NewError(CodeSessionExpired, http.StatusUnauthorized, "会话已过期")
)

// 预定义错误 - 业务逻辑
var (
    ErrBusinessLogic      = NewError(CodeBusinessLogic, http.StatusBadRequest, "业务逻辑错误")
    ErrWorkflowError      = NewError(CodeWorkflowError, http.StatusBadRequest, "工作流错误")
    ErrStatusConflict     = NewError(CodeStatusConflict, http.StatusConflict, "状态冲突")
    ErrOperationFailed    = NewError(CodeOperationFailed, http.StatusInternalServerError, "操作失败")
    ErrResourceConflict   = NewError(CodeResourceConflict, http.StatusConflict, "资源冲突")
    ErrPreconditionFailed = NewError(CodePreconditionFailed, http.StatusPreconditionFailed, "前置条件失败")
    ErrQuotaExceeded      = NewError(CodeQuotaExceeded, http.StatusForbidden, "配额超限")
    ErrResourceExhausted  = NewError(CodeResourceExhausted, http.StatusTooManyRequests, "资源耗尽")
)

// 预定义错误 - 外部服务
var (
    ErrExternalService    = NewError(CodeExternalService, http.StatusBadGateway, "外部服务错误")
    ErrServiceUnavailable = NewError(CodeServiceUnavailable, http.StatusServiceUnavailable, "服务不可用")
    ErrServiceTimeout     = NewError(CodeServiceTimeout, http.StatusRequestTimeout, "服务超时")
    ErrThirdPartyError    = NewError(CodeThirdPartyError, http.StatusBadGateway, "第三方服务错误")
    ErrNetworkError       = NewError(CodeNetworkError, http.StatusBadGateway, "网络错误")
    ErrDatabaseError      = NewError(CodeDatabaseError, http.StatusInternalServerError, "数据库错误")
    ErrCacheError         = NewError(CodeCacheError, http.StatusInternalServerError, "缓存错误")
    ErrMessageQueueError  = NewError(CodeMessageQueueError, http.StatusInternalServerError, "消息队列错误")
)

// 预定义错误 - 系统错误
var (
    ErrInternalError      = NewError(CodeInternalError, http.StatusInternalServerError, "内部错误")
    ErrConfigurationError = NewError(CodeConfigurationError, http.StatusInternalServerError, "配置错误")
    ErrFileSystemError    = NewError(CodeFileSystemError, http.StatusInternalServerError, "文件系统错误")
    ErrMemoryError        = NewError(CodeMemoryError, http.StatusInternalServerError, "内存错误")
    ErrConcurrencyError   = NewError(CodeConcurrencyError, http.StatusInternalServerError, "并发错误")
    ErrDeadlockError      = NewError(CodeDeadlockError, http.StatusInternalServerError, "死锁错误")
)

// 预定义错误 - 限流
var (
    ErrRateLimitExceeded = NewError(CodeRateLimitExceeded, http.StatusTooManyRequests, "请求频率超限")
    ErrTooManyRequests   = NewError(CodeTooManyRequests, http.StatusTooManyRequests, "请求过多")
    ErrConcurrentLimit   = NewError(CodeConcurrentLimit, http.StatusTooManyRequests, "并发数超限")
    ErrAPIQuotaExceeded  = NewError(CodeAPIQuotaExceeded, http.StatusTooManyRequests, "API配额超限")
)

// 预定义错误 - 文件处理
var (
    ErrFileNotFound    = NewError(CodeFileNotFound, http.StatusNotFound, "文件不存在")
    ErrFileTooBig      = NewError(CodeFileTooBig, http.StatusRequestEntityTooLarge, "文件过大")
    ErrInvalidFileType = NewError(CodeInvalidFileType, http.StatusBadRequest, "文件类型无效")
    ErrFileCorrupted   = NewError(CodeFileCorrupted, http.StatusBadRequest, "文件损坏")
    ErrUploadFailed    = NewError(CodeUploadFailed, http.StatusInternalServerError, "上传失败")
    ErrDownloadFailed  = NewError(CodeDownloadFailed, http.StatusInternalServerError, "下载失败")
    ErrFilePermission  = NewError(CodeFilePermission, http.StatusForbidden, "文件权限不足")
)

// 预定义错误 - 安全相关
var (
    ErrEncryptionError    = NewError(CodeEncryptionError, http.StatusInternalServerError, "加密错误")
    ErrDecryptionError    = NewError(CodeDecryptionError, http.StatusInternalServerError, "解密错误")
    ErrSignatureInvalid   = NewError(CodeSignatureInvalid, http.StatusUnauthorized, "签名无效")
    ErrCertificateInvalid = NewError(CodeCertificateInvalid, http.StatusUnauthorized, "证书无效")
    ErrSecurityViolation  = NewError(CodeSecurityViolation, http.StatusForbidden, "安全违规")
)

