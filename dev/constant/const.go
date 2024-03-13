package constant

type RequestIDKey struct{}
type ClientIDKey struct{}

const (
	ReqIDLog = "REQUEST_ID"
)

type ResCode string

const (
	APISuccess       ResCode = "100"
	APIUnauthorized  ResCode = "401"
	APIInternalError ResCode = "500"
)

const (
	InvalidRequest    ResCode = "101"
	InvalidRowData    ResCode = "102"
	EmptyRowMandatory ResCode = "103"
	FileSizeLimit     ResCode = "104"
)