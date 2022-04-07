package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var insertedId int

func TestMain(m *testing.M) {
	ginRouter, pool, _ := SetupRouter()
	router = ginRouter
	var code int = m.Run()
	pool.Close()
	os.Exit(code)
}

func validateJSON(str string) error {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js)
}

func TestGET(t *testing.T) {
	var w *httptest.ResponseRecorder = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/get", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
}

func TestPOST(t *testing.T) {
	var mcPostBody map[string]interface{} = map[string]interface{}{
		"firstname": "User1",
		"lastname":  "User2",
	}
	var w *httptest.ResponseRecorder = httptest.NewRecorder()
	body, _ := json.Marshal(mcPostBody)
	req, _ := http.NewRequest("POST", "/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var user map[string]interface{} = map[string]interface{}{}
	json.Unmarshal([]byte(w.Body.Bytes()), &user)
	insertedId = user["id"].(int)
	// fmt.Println(insertedId)
}

// func TestUpdate(t *testing.T) {
// 	var mcPostBody map[string]interface{} = map[string]interface{}{
// 		"firstname": "User2", // User1 -> User2
// 	}
// 	router, pool, _ := SetupRouter()
// 	defer pool.Close()
// 	var w *httptest.ResponseRecorder = httptest.NewRecorder()
// 	body, _ := json.Marshal(mcPostBody)
// 	req, _ := http.NewRequest("PATCH", "/update", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Code)
// }
