package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateEntry(t *testing.T) {

	var jsonStr = []byte(`{"ProduceCode":"upup-down-left-righ", "UnitPrice": "3.90", "Name": "name"}`)

	req, err := http.NewRequest("POST", "/addItem", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddToServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected := `{"ProduceCode":"upup-down-left-righ", "UnitPrice": "3.90", "Name": "name"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	jsonStr = []byte(`{"ProduceCode":"1234-5678-1234-5678", "UnitPrice": "3.90", "Name": "name"}`)
	req, err = http.NewRequest("POST", "/addItem", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(AddToServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"ProduceCode":"1234-5678-1234-5678", "UnitPrice": "3.90", "Name": "name"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	jsonStr = []byte(`{"ProduceCode":"jan-ken-ro", "UnitPrice": "3.90", "Name": "name"}`)
	req, err = http.NewRequest("POST", "/addItem", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(AddToServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = `{"ProduceCode":"jan-ken-ro", "UnitPrice": "3.90", "Name": "name"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetEntry(t *testing.T) {

	var jsonStr = []byte(``)
	req, err := http.NewRequest("GET", "/item/upup-down-left-righ", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	//Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"Id": "upup-down-left-righ",
	}

	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(GetOneFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var testItem DataStruct
	json.Unmarshal(rr.Body.Bytes(), &testItem)
	var expectedItem DataStruct
	var expectedStr = []byte(`{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":"3.90"}`)
	json.Unmarshal(expectedStr, &expectedItem)
	if testItem != expectedItem {
		t.Errorf("handler returned unexpected body: got %v want %v",
			testItem, expectedItem)
	}

	req, err = http.NewRequest("GET", "/item/450", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	vars = map[string]string{
		"Id": "450",
	}
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(GetOneFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	json.Unmarshal(rr.Body.Bytes(), &testItem)
	expectedStr = []byte(`{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":"3.90"}`)
	json.Unmarshal(expectedStr, &expectedItem)
	if testItem != expectedItem {
		t.Errorf("handler returned unexpected body: got %v want %v",
			testItem, expectedItem)
	}
}
func TestGetAllEntries(t *testing.T) {

	var jsonStr = []byte(``)
	req, err := http.NewRequest("GET", "/items", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var testItem DataStruct
	json.Unmarshal(rr.Body.Bytes(), &testItem)
	var expectedItem DataStruct
	var expectedStr = []byte(`[{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":"3.90"},{"ProduceCode":"1234-5678-1234-5678","Name":"name","UnitPrice":"3.90"}]`)
	json.Unmarshal(expectedStr, &expectedItem)
	if testItem != expectedItem {
		t.Errorf("handler returned unexpected body: got %v want %v",
			testItem, expectedItem)
	}
}

// tests one correct, one incorrected, and one with wrong type entirely
func TestDeleteEntries(t *testing.T) {
	var jsonStr = []byte(`{"ProduceCode":"upup-down-left-righ", "UnitPrice": "3.90", "Name": "name"}`)

	req, err := http.NewRequest("DELETE", "/delete/upup-down-left-righ", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"Id": "upup-down-left-righ",
	}
	req = mux.SetURLVars(req, vars)
	handler := http.HandlerFunc(DeleteFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expectedStr := []byte(`{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":"3.90"}`)
	var testItem DataStruct
	json.Unmarshal(rr.Body.Bytes(), &testItem)
	var expectedItem DataStruct
	json.Unmarshal(expectedStr, &expectedItem)
	if testItem != expectedItem {
		t.Errorf("handler returned unexpected body: got %v want %v",
			testItem, expectedItem)
	}

	jsonStr = []byte(`{"ProduceCode":"upup-down-left-righ", "UnitPrice": "3.90", "Name": "name"}`)
	req, err = http.NewRequest("Delete", "/delete/what-is-the-greatest-number-ever", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	vars = map[string]string{
		"Id": "what-is-the-greatest-number-ever",
	}
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(DeleteFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `Wrong Method Delete used `
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	req, err = http.NewRequest("POST", "/delete/what-is-the-greatest-number-ever", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	vars = map[string]string{
		"Id": "what-is-the-greatest-number-ever",
	}
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(DeleteFromServer)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	json.Unmarshal(rr.Body.Bytes(), &testItem)
	expectedStr = []byte(`[{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":"3.90"},{"ProduceCode":"1234-5678-1234-5678","Name":"name","UnitPrice":"3.90"}]`)
	expected = `Wrong Method POST used `
	json.Unmarshal(expectedStr, &expectedItem)

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
