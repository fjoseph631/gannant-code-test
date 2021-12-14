package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type DataStruct struct {
	ProduceCode string `json:"ProduceCode"`
	Name        string `json:"Name"`
	UnitPrice   string `json:"UnitPrice"`
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
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(ProduceStore)
}

func GetOneFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, "Wrong Method %v used", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	key := params["Id"]
	for _, produceItem := range ProduceStore {
		if produceItem.ProduceCode == key {
			json.NewEncoder(w).Encode(produceItem)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func DeleteFromServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodGet {
		fmt.Fprintf(w, "Wrong Method %v used ", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	key := params["Id"]
	for index, produceItem := range ProduceStore {
		// if our id path parameter matches one of our
		// produceItems
		if produceItem.ProduceCode == key {
			// updates our produceItems array to remove the produceItem
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(produceItem)
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
	w.WriteHeader(http.StatusBadRequest)
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
	var produceItem DataStruct
	json.Unmarshal(reqBody, &produceItem)
	// regex to match desired product code
	re, _ := regexp.Compile("[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]")
	// this would allow for one digit after place, but couldn't find a regex to achieve that
	re2, _ := regexp.Compile(`\d+\.\d{2}`)
	match := re.MatchString(produceItem.ProduceCode)
	match2 := re2.MatchString(produceItem.UnitPrice)
	if match == false {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if match2 == false {
			fmt.Println("Unit Price requires two decimal, and must be a float")
			log.Println("unit price code wrong", match2)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			if produceItem.Name != "" && produceItem.UnitPrice != "" {
				// update our global Articles array to include
				// our new Article
				found := false
				for i := 0; i < len(ProduceStore); i++ {
					if ProduceStore[i].ProduceCode == produceItem.ProduceCode {
						ProduceStore[i].Name = produceItem.Name
						ProduceStore[i].UnitPrice = produceItem.UnitPrice
						found = true
						break
					}

				}
				if !found {
					ProduceStore = append(ProduceStore, produceItem)
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}