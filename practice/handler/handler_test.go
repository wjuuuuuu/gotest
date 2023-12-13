package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleSendJSON() {
	r := httptest.NewRequest("GET", "/sendjson", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)

	var u struct {
		Name  string
		Email string
	}
	if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
		log.Println("ERROR: ", err)
	}
	fmt.Println(u)
}
