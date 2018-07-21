// Copyright (c) 2018 Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	var address = flag.String("address", ":8080", "TCP address to bind to.")
	var socket = flag.String("socket", "", "UNIX socket to bind to.")
	flag.Parse()

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/ping", handlePing)
	var err error

	if *socket != "" {
		// If it isn't there, the error is harmless. If it can't be removed,
		// we'll get the error again on the next line anyway.
		_ = os.Remove(*socket)
		var l net.Listener
		l, err = net.Listen("unix", *socket)
		if err != nil {
			log.Fatal(err)
		}
		// Clean up on close.
		defer os.Remove(*socket)
		defer l.Close()
		// Make sure the webserver can read it.
		err = os.Chmod(*socket, 0777)
		if err != nil {
			log.Fatal(err)
		}
		// Serve HTTP over the socket.
		err = http.Serve(l, nil)
	} else {
		err = http.ListenAndServe(*address, nil)
	}
	log.Fatal("Exiting: ", err)
}

func handlePing(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusMethodNotAllowed)
	return
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
