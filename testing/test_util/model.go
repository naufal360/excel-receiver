package test_util

type MultipartRequest struct {
	Field string
	Value string
}

type JSONRequest struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	Src  string `json:"src"`
}
