package errors

import "errors"

var (
	NotSpecified        = errors.New("NotSpecified")
	OpenError           = errors.New("OpenError")
	CloseError          = errors.New("CloseError")
	DeviceError         = errors.New("DeviceError")
	DeadlockDetected    = errors.New("DeadlockDetected")
	InvalidArgument     = errors.New("InvalidArgument")
	InvalidData         = errors.New("InvalidData")
	InvalidOperation    = errors.New("InvalidOperation")
	InvalidResult       = errors.New("InvalidResult")
	InvalidType         = errors.New("InvalidType")
	MetadataError       = errors.New("MetadataError")
	NotImplemented      = errors.New("NotImplemented")
	NotInstalled        = errors.New("NotInstalled")
	ObjectNotFound      = errors.New("ObjectNotFound")
	OperationStopped    = errors.New("OperationStopped")
	OperationTimeout    = errors.New("OperationTimeout")
	SyntaxError         = errors.New("SyntaxError")
	ParserError         = errors.New("ParserError")
	PermissionDenied    = errors.New("PermissionDenied")
	ResourceBusy        = errors.New("ResourceBusy")
	ResourceExists      = errors.New("ResourceExists")
	ResourceUnavailable = errors.New("ResourceUnavailable")
	ReadError           = errors.New("ReadError")
	WriteError          = errors.New("WriteError")
	FromStdErr          = errors.New("FromStdErr")
	SecurityError       = errors.New("SecurityError")
	ProtocolError       = errors.New("ProtocolError")
	ConnectionError     = errors.New("ConnectionError")
	AuthenticationError = errors.New("AuthenticationError")
	LimitsExceeded      = errors.New("LimitsExceeded")
	QuotaExceeded       = errors.New("QuotaExceeded")
	NotEnabled          = errors.New("NotEnabled")
)

// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.errorcategory?view=powershellsdk-7.4.0
var categoryEnum = map[int]error{
	0:  NotSpecified,
	1:  OpenError,
	2:  CloseError,
	3:  DeviceError,
	4:  DeadlockDetected,
	5:  InvalidArgument,
	6:  InvalidData,
	7:  InvalidOperation,
	8:  InvalidResult,
	9:  InvalidType,
	10: MetadataError,
	11: NotImplemented,
	12: NotInstalled,
	13: ObjectNotFound,
	14: OperationStopped,
	15: OperationTimeout,
	16: SyntaxError,
	17: ParserError,
	18: PermissionDenied,
	19: ResourceBusy,
	20: ResourceExists,
	21: ResourceUnavailable,
	22: ReadError,
	23: WriteError,
	24: FromStdErr,
	25: SecurityError,
	26: ProtocolError,
	27: ConnectionError,
	28: AuthenticationError,
	29: LimitsExceeded,
	30: QuotaExceeded,
	31: NotEnabled,
}

func GetCategoryEnumItem(code int) error {
	return categoryEnum[code]
}
