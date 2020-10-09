package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

type Response struct {
	Location string `json:"location"`
	Value    string `json:"value"`
}

type Body struct {
	Value string `json:"value"`
}

var node kademlia.Node

func API(output io.Writer, n kademlia.Node) {
	fmt.Println("Starting REST API")
	node = n

	r := mux.NewRouter()
	r.HandleFunc("/objects/{hash}", GetHandler).Methods("GET")
	r.HandleFunc("/objects", PostHandler).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", r))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.URL.Path, "/")[2]

	if len(hash) != 40 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		value, err := node.FindValue(hash)

		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusNotFound)
		} else {
			res := Response{"/objects/" + hash, value}
			json.NewEncoder(w).Encode(res)
		}
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	body := Body{}
	err = json.Unmarshal(b, &body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		key := node.StoreValue(body.Value)
		res := Response{"/objects/" + key, body.Value}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}
