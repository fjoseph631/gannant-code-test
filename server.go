package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

)

type DataStruct struct {
	ProduceCode string  `json:"ProduceCode"`
	Name        string  `json:"Name"`
	UnitPrice   float64 `json:"UnitPrice"`
}
type IdStruct struct {
	Id string  `json:"Id"`
}
//var ProduceStore []DataStruct
var ProduceStore = make(map[string]DataStruct)
var re, errRe = regexp.Compile("[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]")

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloServer)
	mux.HandleFunc("/items", HandleRequests)
	fmt.Println("Server started at port 8080")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Server")
}
func HandleRequests(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("content-type") != "application/json"{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Header Type %v used", r.Header.Values("content-type"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetAllFromServer(w,r);
		break;
	case http.MethodPost:
		AddToServer(w,r);
		break;
	case http.MethodDelete:
		DeleteFromServer(w,r);
		break;
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		break;
	}
}

func GetAllFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		return
	}
		Id := r.URL.Query().Get("Id");
	if (Id != ``){
		produceItem, exists := ProduceStore[strings.ToLower(Id)]
			if exists{
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("Product Code: " + produceItem.ProduceCode)
				json.NewEncoder(w).Encode("UnitPrice: " +  strconv.FormatFloat(produceItem.UnitPrice, 'f', 5, 64))
				json.NewEncoder(w).Encode("Name: " + produceItem.Name)
			} else {
				fmt.Fprintf(w, "%s Not found", Id)
				w.WriteHeader(http.StatusNotFound)
			}
			return;
	}
	for _, v := range ProduceStore {
		json.NewEncoder(w).Encode(v)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used ", r.Method)
		return
	}
	Id := r.URL.Query().Get("Id");
	// if our id path parameter matches one of our
	// produceItems
	_, exists := ProduceStore[strings.ToLower(Id)]
	if exists{
		fmt.Fprintf(w, "Deleted %s", Id);
		w.WriteHeader(202)
		delete(ProduceStore, Id) 
	} else
	{
		fmt.Fprintf(w, "%s not found", Id)
		w.WriteHeader(http.StatusNotFound)
	}
	
}

func AddToServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Read all failed %v", err)
		return;
	}
	var produceItem DataStruct
	err = json.Unmarshal(reqBody, &produceItem)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Json unmarshal failed %v", err)
		return;
	}
	strconv.FormatFloat(produceItem.UnitPrice, 'f', 2, 32)
	// regex to match desired product code
	if errRe != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Read all failed %s", err)
		return;
	}
	match := re.MatchString(produceItem.ProduceCode)
	if match == false {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Produce Code %s does not match required regex", produceItem.ProduceCode);
		return;
	} else {
		// Default Price of 0.0 should be rejected
		if produceItem.Name != "" && produceItem.UnitPrice > 0.0 {
			// update global produce store
			_, exists := ProduceStore[strings.ToLower(produceItem.ProduceCode)]
			if exists{
				fmt.Fprintf(w, "Already exists in store, updating with new value");
			} else
			{
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, "Added %s to store", produceItem.ProduceCode)
			}
			ProduceStore[strings.ToLower(produceItem.ProduceCode)] = produceItem;
		}else
		{
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Item was invalid");
		}
	}

}
