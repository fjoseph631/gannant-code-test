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

	"github.com/gorilla/mux"
)

type DataStruct struct {
	ProduceCode string  `json:"ProduceCode"`
	Name        string  `json:"Name"`
	UnitPrice   float64 `json:"UnitPrice"`
}

var ProduceStore []DataStruct

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", HelloServer)
	router.HandleFunc("/addItem", AddToServer)
	router.HandleFunc("/items", GetAllFromServer)
	router.HandleFunc("/item/{Id}", GetOneFromServer)
	router.HandleFunc("/delete/{Id}", DeleteFromServer).Methods(http.MethodDelete, http.MethodGet)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Server")
}

func GetAllFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		return
	}
	json.NewEncoder(w).Encode(ProduceStore)
}

func GetOneFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		return
	}
	params := mux.Vars(r)
	key := params["Id"]
	for _, produceItem := range ProduceStore {
		if strings.ToLower(produceItem.ProduceCode) == strings.ToLower(key) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(produceItem)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func DeleteFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong Method %v used ", r.Method)
		return
	}
	params := mux.Vars(r)
	key := params["Id"]
	for index, produceItem := range ProduceStore {
		// if our id path parameter matches one of our
		// produceItems
		if strings.ToLower(produceItem.ProduceCode) == strings.ToLower(key) {
			// updates our produceItems array to remove the produceItem
			w.WriteHeader(http.StatusOK)
			ProduceStore = append(ProduceStore[:index], ProduceStore[index+1:]...)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func AddToServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var produceItem DataStruct
	json.Unmarshal(reqBody, &produceItem)
	strconv.FormatFloat(produceItem.UnitPrice, 'f', 2, 32)
	// regex to match desired product code
	re, _ := regexp.Compile("[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]")
	match := re.MatchString(produceItem.ProduceCode)
	if match == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Default Price of 0.0 should be rejected
		if produceItem.Name != "" && produceItem.UnitPrice > 0.0 {
			// update global produce store
			found := false
			for i := 0; i < len(ProduceStore); i++ {
				if strings.ToLower(ProduceStore[i].ProduceCode) == strings.ToLower(produceItem.ProduceCode) {
					ProduceStore[i].Name = produceItem.Name
					ProduceStore[i].UnitPrice = produceItem.UnitPrice
					found = true
					break
				}

			}
			if !found {
				ProduceStore = append(ProduceStore, produceItem)
				w.WriteHeader(http.StatusCreated)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
