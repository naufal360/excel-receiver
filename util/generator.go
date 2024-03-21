package util

import "github.com/thanhpk/randstr"

func GenerateReqID() string {
	return randstr.String(6)
}
