package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var insertedId int32

func TestMain(m *testing.M) {
	ginRouter, pool, _ := SetupRouter()
	router = ginRouter
	insertedRow, _ := pool.Query(context.Background(), "INSERT INTO mytable (firstname, lastname) VALUES ('User1', 'User2') RETURNING id")
	var insertedUser map[string]interface{} = map[string]interface{}{}
	if err := pgxscan.ScanOne(&insertedUser, insertedRow); err != nil {
		log.Fatal(err)
	}
	insertedId = insertedUser["id"].(int32)
	var code int = m.Run()
	pool.Close()
	os.Exit(code)
}

func validateJSON(str string) error {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js)
}

func containsMap(mapA map[string]interface{}, mapB map[string]interface{}) bool {
	for key := range mapB {
		if mapA[key] != mapB[key] {
			return false
		}
	}
	return true
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
	var insertedUser map[string]interface{} = map[string]interface{}{}
	json.Unmarshal(w.Body.Bytes(), &insertedUser)
	var id float64 = insertedUser["id"].(float64)
	// Check that user was created
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/get", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var users []map[string]interface{} = make([]map[string]interface{}, 0)
	json.Unmarshal(w.Body.Bytes(), &users)
	var insertedSuccessfully bool = false
	for _, user := range users {
		if user["id"] == id {
			if containsMap(user, mcPostBody) {
				insertedSuccessfully = true
			}
		}
	}
	if !insertedSuccessfully {
		t.Fatalf("Not inserted successfully")
	}
	// Delete user
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/delete/"+fmt.Sprint(id), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var deletedUser map[string]interface{} = map[string]interface{}{}
	json.Unmarshal(w.Body.Bytes(), &deletedUser)
	var eq bool = reflect.DeepEqual(deletedUser, insertedUser)
	if !eq {
		t.Fatalf("Inserted user and deleted user are not the same")
	}
}

func TestUpdate(t *testing.T) {
	var mcPostBody map[string]interface{} = map[string]interface{}{
		"firstname": "User2", // User1 -> User2
	}
	var w *httptest.ResponseRecorder = httptest.NewRecorder()
	body, _ := json.Marshal(mcPostBody)
	req, _ := http.NewRequest("PATCH", "/update/"+fmt.Sprint(insertedId), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var updatedUser map[string]interface{} = map[string]interface{}{}
	json.Unmarshal(w.Body.Bytes(), &updatedUser)
	var id float64 = updatedUser["id"].(float64)
	// Check user was updated
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/get", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var users []map[string]interface{} = make([]map[string]interface{}, 0)
	json.Unmarshal(w.Body.Bytes(), &users)
	var updatedSuccessfully bool = false
	for _, user := range users {
		if user["id"] == id {
			if containsMap(user, mcPostBody) {
				updatedSuccessfully = true
			}
		}
	}
	if !updatedSuccessfully {
		t.Fatalf("Not inserted successfully")
	}
}
