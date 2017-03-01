// +build integration
package main

import (
	"testing"
	"./model"
	"net/http/httptest"
	"net/http"
	"strings"
	"time"
	"net/url"
	"log"
)

type mockDB struct {}

var dataTests = []struct {
	testName	string
	request 	string
	expectedStatus 	int
}{
	{
		"Valid request",
		"{\"nodeId\":1, \"timestamp\":\""+time.Now().Format(time.RFC3339)+"\", \"data\":{\"temp\":\"2.5\"}}",
		http.StatusOK,
	},
	{
		"Empty data object",
		"{\"nodeId\":1, \"timestamp\":\""+time.Now().Format(time.RFC3339)+"\", \"data\":{}}",
		http.StatusBadRequest,
	},
	{
		"Missing timestamp",
		"{\"nodeId\":1, \"data\":{}}",
		http.StatusBadRequest,
	},
	{
		"Array instead of object for data array",
		"{\"nodeId\":1, \"data\":[]}",
		http.StatusBadRequest,
	},
	{
		"Malformed JSON",
		"{\"1\",\"",
		http.StatusBadRequest,
	},
}

var dataRequestTests = []struct {
	testName	string
	queryString	string
	expectedStatus 	int
}{
	{
		"Valid request",
		"nodeId=1&start="+time.Now().Format(time.RFC3339)+"&end="+time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		http.StatusOK,
	},
	{
		"End before start time",
		"nodeId=1&start="+time.Now().Format(time.RFC3339)+"&end="+time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		http.StatusBadRequest,
	},/*
	{
		"Missing node id",
		"{\"start\":\""+time.Now().Format(time.RFC3339)+"\", \"end\":\""+time.Now().Add(1 * time.Hour).Format(time.RFC3339)+"\"}",
		http.StatusBadRequest,
	},
	{
		"Missing start",
		"{\"nodeId\":1, \"end\":\""+time.Now().Add(1 * time.Hour).Format(time.RFC3339)+"\"}",
		http.StatusBadRequest,
	},
	{
		"Missing end",
		"{\"nodeId\":1, \"start\":\""+time.Now().Format(time.RFC3339)+"\"}",
		http.StatusBadRequest,
	},*/
}

func (mdb *mockDB) GetLast(int)(*model.Data, error) {
	return nil, nil
}

func (mdb *mockDB) SaveData(data *model.Data) (error) {
	return nil
}

func (mdb *mockDB) GetData(request *model.DataRequest) ([]map[string]interface{}, error) {
	return nil, nil
}

func TestHandlePostData(t *testing.T) {
	for _, tt := range dataTests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/data", strings.NewReader(tt.request))

		env := Env{db: &mockDB{}}

		http.HandlerFunc(env.handleDataPost).ServeHTTP(rec, req)

		if rec.Code != tt.expectedStatus {
			t.Errorf("[!!] %s: Expected %d got %d", tt.testName, tt.expectedStatus, rec.Code)
		} else {
			t.Log("[OK]", tt.testName)
		}
	}
}

func TestHandleValidate(t *testing.T) {
	for _, tt := range dataTests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/data", strings.NewReader(tt.request))

		env := Env{db: &mockDB{}}

		http.HandlerFunc(env.handleDataValidate).ServeHTTP(rec, req)

		if rec.Code != tt.expectedStatus {
			t.Errorf("[!!] %s: Expected %d got %d", tt.testName, tt.expectedStatus, rec.Code)
		} else {
			t.Log("[OK]", tt.testName)
		}
	}
}

func TestHandleGetData(t *testing.T) {
	for _, tt := range dataRequestTests {
		parsedQueryStr, _ := url.ParseQuery(tt.queryString)

		log.Println(parsedQueryStr)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/data", nil)

		req.URL.RawQuery = parsedQueryStr.Encode()

		env := Env{db: &mockDB{}}

		http.HandlerFunc(env.handleDataGet).ServeHTTP(rec, req)



		if rec.Code != tt.expectedStatus {
			t.Errorf("[!!] %s: Expected %d got %d", tt.testName, tt.expectedStatus, rec.Code)
		} else {
			t.Log("[OK]", tt.testName)
		}
	}
}