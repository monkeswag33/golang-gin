package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/monkeswag33/golang-gin/global"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var router *gin.Engine
var insertedUser global.User
var db *gorm.DB

func TestMain(m *testing.M) {
	ginRouter, localDB, _ := SetupRouter()
	router = ginRouter
	db = localDB
	insertedUser = global.User{
		Firstname: "User1",
		Lastname:  "User2",
	}
	db.Create(&insertedUser)
	var code int = m.Run()
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
	// Verify that it is correct
	var requestedUsers []global.User
	if err := json.Unmarshal(w.Body.Bytes(), &requestedUsers); err != nil {
		t.Fatal(err)
	}
	var users []global.User
	db.Find(&users)
	if !reflect.DeepEqual(requestedUsers, users) {
		t.Fatalf("Requested users and users retreived from database are not the same")
	}
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
	var insertedUser global.User
	json.Unmarshal(w.Body.Bytes(), &insertedUser)
	// Check that user was created
	var dbUser global.User
	db.Find(&dbUser, insertedUser.ID)
	if !reflect.DeepEqual(insertedUser, dbUser) {
		t.Fatalf("Inserted user and requested user are not the same")
	}
	// Delete user
	db.Delete(&dbUser)
}

func TestUpdate(t *testing.T) {
	var mcPostBody map[string]interface{} = map[string]interface{}{
		"firstname": "User2", // User1 -> User2
	}
	var w *httptest.ResponseRecorder = httptest.NewRecorder()
	body, _ := json.Marshal(mcPostBody)
	req, _ := http.NewRequest("PATCH", "/update/"+fmt.Sprint(insertedUser.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var updatedUser global.User
	json.Unmarshal(w.Body.Bytes(), &updatedUser)
	// Check if user was updated
	var dbUser global.User
	db.Find(&dbUser, insertedUser.ID)
	if !reflect.DeepEqual(dbUser, updatedUser) {
		t.Fatalf("Not updated successfully")
	}
}

func TestDelete(t *testing.T) {
	var w *httptest.ResponseRecorder = httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete/"+fmt.Sprint(insertedUser.ID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, validateJSON(w.Body.String()))
	var deletedUser global.User
	json.Unmarshal(w.Body.Bytes(), &deletedUser)
	// Check that user was deleted
	var exists bool
	db.Model(&global.User{}).Select("COUNT(*) > 0").Find(&exists, insertedUser.ID)
	if exists {
		t.Fatalf("Record was not deleted")
	}
}
