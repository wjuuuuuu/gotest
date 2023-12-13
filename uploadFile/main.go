package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/uploads", uploadHandler) // 이 경로로 들어오면 uploadHandler 함수를 탄다
	http.Handle("/", http.FileServer(http.Dir("uploadFile/public")))
	err := http.ListenAndServe(":1324", nil)
	if err != nil {
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upload_file") // html form 의 file type, name이 upload_file로 들어온 request
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	dirname := "C:/goproject/uploadFile/uploads"
	os.Mkdir(dirname, 0777) // 경로 생성
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath) // 경로에 파일 생성
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	io.Copy(file, uploadFile) // io.copy(writer, reader)  reader로 읽어서 writer에 복사
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)

	//test
}
