package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)
type baseTest struct {
	name string
	input string
	id string
	expected string
	requestType string
	status int
	header string
	queryKey string
	queryValue string
    data DataStruct

}
var tests = []baseTest {
	{
		name:"correctly built add item request", requestType:"POST", header:"application/json",
		expected: `Added upup-down-left-righ to store`,
		input: `{"ProduceCode":"upup-down-left-righ", "UnitPrice": 3.90, "Name": "name"}`,
		status: http.StatusCreated,
		data: DataStruct {
	
		},	
	},
	{
		name:"correctly built get request", requestType:"GET", header:"application/json",
		expected: `{"ProduceCode":"upup-down-left-righ","Name":"name","UnitPrice":3.9}`,
		input: ``,
		queryKey: `Id`,
		queryValue: `upup-downd-left-righ`,
		status: http.StatusOK,
		data: DataStruct {
		},	
	},
	{
		name:"correctly built get request but single item and not found", requestType:"GET", header:"application/json",
		expected: `upup-downd-left-righ Not found`,
		input: ``,
		queryKey: ``,
		queryValue: ``,
		status: http.StatusOK,
		data: DataStruct{
		},	
	},
	{
		name:"correctly built delete request", requestType:"DELETE", header:"application/json",
		input: ``,
		expected: `Deleted upup-down-down-left-righ`,
		queryKey: `Id`,
		queryValue: `upup-down-left-righ`,
		status: http.StatusOK,
		data: DataStruct {
			ProduceCode: "upup-down-left-righ",
			UnitPrice: 3.90,
			Name: "name",
		},	
	},
	{
		name:"wrong header used add item", requestType:"POST", header:"plaintext",
		expected: `Wrong Header Type [plaintext] used`,
		input: `{"ProduceCode":"upup-down-left-righ", "UnitPrice": 3.90, "Name": "name"}`,
		status: http.StatusBadRequest,
		data: DataStruct {
			ProduceCode: ``,
			UnitPrice: 3.90,
			Name: ``,
		},	
	},
	{
		name:"regex failed", requestType:"POST", header:"application/json",
		input: `{"ProduceCode":"upup-downhkjh-left-righ", "UnitPrice": 3.90, "Name": "name"}`,
		expected: "Produce Code upup-downhkjh-left-righ does not match required regex",
		status: http.StatusBadRequest,
		data: DataStruct {
			ProduceCode: "",
			UnitPrice: 0.0,
			Name: "",
		},	
	},
	{
		name:"json unmarshal failed", requestType:"POST", header:"application/json",
		input: `{ProduceCode":"upup-down-left-righ", "UnitPrice": 3.90, "Nae": "name"}`,
		expected: `Deleted upup-down-down-left-righ`,
		status: http.StatusBadRequest,
		data: DataStruct {
			ProduceCode: "upup-down-left-righ",
			UnitPrice: 3.90,
			Name: "name",
		},	
	}}
func TestGet(t *testing.T) {
    
	for _, tc := range tests {
		var jsonStr = []byte(tc.input)
		req, _ := http.NewRequest(tc.requestType,"/items/"+tc.id, bytes.NewBuffer(jsonStr));
		req.Header.Set("Content-Type", tc.header)
		q := req.URL.Query()
		q.Add(tc.queryKey, tc.queryValue)
		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleRequests)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != tc.status {
			t.Errorf("handler returned wrong status code: got %v want %v case %v output %v",
				status, tc.status, tc.name, rr.Body )
		}
	}
}
