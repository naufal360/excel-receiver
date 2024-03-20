package test_util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"testing"
)

func CreateMultipartFormData(t *testing.T, fieldName, fileName string, multipartData MultipartRequest) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	// for _, val := range otherData {
	err = w.WriteField(multipartData.Field, multipartData.Value)
	if err != nil {
		t.Errorf("Error creating writer field: %v", err)
	}
	// }
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}
