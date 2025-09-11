package errorx

// ErrorCode 错误码类型
type ErrorCode int

const (
    // 1000-1099: 数据相关错误
    CodeRecordNotFound     ErrorCode = 1001
    CodeRecordDuplicated   ErrorCode = 1002
    CodeDataCorrupted      ErrorCode = 1003
    CodeDataTooLarge       ErrorCode = 1004
    CodeDataValidationFail ErrorCode = 1005
    CodeConstraintViolated ErrorCode = 1006
    CodeDataExpired        ErrorCode = 1007
    CodeDataLocked         ErrorCode = 1008

    // 1100-1199: 请求相关错误
    CodeBadRequest        ErrorCode = 1101
    CodeMissingParameter  ErrorCode = 1102
    CodeInvalidParameter  ErrorCode = 1103
    CodeParameterTooLong  ErrorCode = 1104
    CodeParameterTooShort ErrorCode = 1105
    CodeInvalidFormat     ErrorCode = 1106
    CodeUnsupportedMethod ErrorCode = 1107
    CodeRequestTooLarge   ErrorCode = 1108
    CodeInvalidJSON       ErrorCode = 1109
    CodeInvalidXML        ErrorCode = 1110

    // 1200-1299: 认证授权错误
    CodeUnauthorized       ErrorCode = 1201
    CodeForbidden          ErrorCode = 1202
    CodeTokenExpired       ErrorCode = 1203
    CodeTokenInvalid       ErrorCode = 1204
    CodeTokenMissing       ErrorCode = 1205
    CodePermissionDenied   ErrorCode = 1206
    CodeAccountDisabled    ErrorCode = 1207
    CodeAccountLocked      ErrorCode = 1208
    CodeInvalidCredentials ErrorCode = 1209
    CodeSessionExpired     ErrorCode = 1210

    // 1300-1399: 业务逻辑错误
    CodeBusinessLogic      ErrorCode = 1301
    CodeWorkflowError      ErrorCode = 1302
    CodeStatusConflict     ErrorCode = 1303
    CodeOperationFailed    ErrorCode = 1304
    CodeResourceConflict   ErrorCode = 1305
    CodePreconditionFailed ErrorCode = 1306
    CodeQuotaExceeded      ErrorCode = 1307
    CodeResourceExhausted  ErrorCode = 1308

    // 1400-1499: 外部服务错误
    CodeExternalService    ErrorCode = 1401
    CodeServiceUnavailable ErrorCode = 1402
    CodeServiceTimeout     ErrorCode = 1403
    CodeThirdPartyError    ErrorCode = 1404
    CodeNetworkError       ErrorCode = 1405
    CodeDatabaseError      ErrorCode = 1406
    CodeCacheError         ErrorCode = 1407
    CodeMessageQueueError  ErrorCode = 1408

    // 1500-1599: 系统错误
    CodeInternalError      ErrorCode = 1501
    CodeConfigurationError ErrorCode = 1502
    CodeFileSystemError    ErrorCode = 1503
    CodeMemoryError        ErrorCode = 1504
    CodeConcurrencyError   ErrorCode = 1505
    CodeDeadlockError      ErrorCode = 1506

    // 1600-1699: 限流和频率控制
    CodeRateLimitExceeded ErrorCode = 1601
    CodeTooManyRequests   ErrorCode = 1602
    CodeConcurrentLimit   ErrorCode = 1603
    CodeAPIQuotaExceeded  ErrorCode = 1604

    // 1700-1799: 文件和上传错误
    CodeFileNotFound    ErrorCode = 1701
    CodeFileTooBig      ErrorCode = 1702
    CodeInvalidFileType ErrorCode = 1703
    CodeFileCorrupted   ErrorCode = 1704
    CodeUploadFailed    ErrorCode = 1705
    CodeDownloadFailed  ErrorCode = 1706
    CodeFilePermission  ErrorCode = 1707

    // 1800-1899: 加密和安全错误
    CodeEncryptionError    ErrorCode = 1801
    CodeDecryptionError    ErrorCode = 1802
    CodeSignatureInvalid   ErrorCode = 1803
    CodeCertificateInvalid ErrorCode = 1804
    CodeSecurityViolation  ErrorCode = 1805
)

