package main

import (
	"encoding/json"
	"net/http"
)

func main() {
}

func test(w http.ResponseWriter, req *http.Request) {
	obj := map[string]interface{}{
		"name":     "gee",
		"password": "gee",
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
