package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

var node kademlia.Node

func API(output io.Writer, n kademlia.Node) {
	fmt.Println("Starting REST API")
	node = n

	r := mux.NewRouter()
	r.HandleFunc("/{key}", GetHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", r))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler executing")

	// remove the '/' in the beginning of the path
	// the rest is the hash inputed by the user
	hash := r.URL.Path[1:]

	if len(hash) != 40 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		value, err := node.FindValue(hash)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(value)
		}
	}
}
