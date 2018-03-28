// Copyright (c) 2018 Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	var stateParam = r.URL.Query().Get("state")
	// This will have a token and other fields, but we only care about the host.
	var state struct{ Host string }
	var err = json.Unmarshal([]byte(stateParam), &state)
	if err != nil || state.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var url = *r.URL
	url.Scheme = "https"
	url.Host = state.Host
	log.Println(r.URL.String(), "->", url.String())
	http.Redirect(w, r, url.String(), http.StatusFound)
}
