package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	path := "C:/goproject/uploadFile/uploads/backend.png"

	file, _ := os.Open(path)
	defer file.Close()

	os.RemoveAll("C:/goproject/uploadFile/uploads")

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	assert.NoError(err)
	io.Copy(multi, file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	uploadHandler(res, req)
	assert.Equal(http.StatusOK, res.Code)

	uplaodFilePath := "C:/goproject/uploadFile/uploads/" + filepath.Base(path)
	_, err = os.Stat(uplaodFilePath) // 경로 안에 파일이 잘 들어 있는지 확인
	fmt.Println(err)
	assert.NoError(err)

	uploadFile, _ := os.Open(uplaodFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}

	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(uploadData, originData)

}
