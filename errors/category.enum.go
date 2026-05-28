package errors

import "errors"

var (
	ErrNotSpecified        = errors.New("NotSpecified")
	ErrOpenError           = errors.New("OpenError")
	ErrCloseError          = errors.New("CloseError")
	ErrDeviceError         = errors.New("DeviceError")
	ErrDeadlockDetected    = errors.New("DeadlockDetected")
	ErrInvalidArgument     = errors.New("InvalidArgument")
	ErrInvalidData         = errors.New("InvalidData")
	ErrInvalidOperation    = errors.New("InvalidOperation")
	ErrInvalidResult       = errors.New("InvalidResult")
	ErrInvalidType         = errors.New("InvalidType")
	ErrMetadataError       = errors.New("MetadataError")
	ErrNotImplemented      = errors.New("NotImplemented")
	ErrNotInstalled        = errors.New("NotInstalled")
	ErrObjectNotFound      = errors.New("ObjectNotFound")
	ErrOperationStopped    = errors.New("OperationStopped")
	ErrOperationTimeout    = errors.New("OperationTimeout")
	ErrSyntaxError         = errors.New("SyntaxError")
	ErrParserError         = errors.New("ParserError")
	ErrPermissionDenied    = errors.New("PermissionDenied")
	ErrResourceBusy        = errors.New("ResourceBusy")
	ErrResourceExists      = errors.New("ResourceExists")
	ErrResourceUnavailable = errors.New("ResourceUnavailable")
	ErrReadError           = errors.New("ReadError")
	ErrWriteError          = errors.New("WriteError")
	ErrFromStdErr          = errors.New("FromStdErr")
	ErrSecurityError       = errors.New("SecurityError")
	ErrProtocolError       = errors.New("ProtocolError")
	ErrConnectionError     = errors.New("ConnectionError")
	ErrAuthenticationError = errors.New("AuthenticationError")
	ErrLimitsExceeded      = errors.New("LimitsExceeded")
	ErrQuotaExceeded       = errors.New("QuotaExceeded")
	ErrNotEnabled          = errors.New("NotEnabled")
)

// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.errorcategory?view=powershellsdk-7.4.0
var categoryEnum = map[int]error{
	0:  ErrNotSpecified,
	1:  ErrOpenError,
	2:  ErrCloseError,
	3:  ErrDeviceError,
	4:  ErrDeadlockDetected,
	5:  ErrInvalidArgument,
	6:  ErrInvalidData,
	7:  ErrInvalidOperation,
	8:  ErrInvalidResult,
	9:  ErrInvalidType,
	10: ErrMetadataError,
	11: ErrNotImplemented,
	12: ErrNotInstalled,
	13: ErrObjectNotFound,
	14: ErrOperationStopped,
	15: ErrOperationTimeout,
	16: ErrSyntaxError,
	17: ErrParserError,
	18: ErrPermissionDenied,
	19: ErrResourceBusy,
	20: ErrResourceExists,
	21: ErrResourceUnavailable,
	22: ErrReadError,
	23: ErrWriteError,
	24: ErrFromStdErr,
	25: ErrSecurityError,
	26: ErrProtocolError,
	27: ErrConnectionError,
	28: ErrAuthenticationError,
	29: ErrLimitsExceeded,
	30: ErrQuotaExceeded,
	31: ErrNotEnabled,
}

func GetCategoryEnumItem(code int) error {
	return categoryEnum[code]
}
