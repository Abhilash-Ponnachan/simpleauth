package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// handler funcs for diff req
type reqHandler struct {
	filehandler http.Handler
}

// handler func for /hello[?name=xxx]
func (rh *reqHandler) hello(w http.ResponseWriter, r *http.Request) {
	q, ok := r.URL.Query()["name"]
	if ok {
		fmt.Fprintf(w, "<h1>Salut, Bonjour  %s!</h1>", string(q[0]))
	} else {
		fmt.Fprint(w, "<h1>Salut, Bonjour!</h1>")
	}
}

// handler func for /datetime => JSON
func (rh *reqHandler) datetime(w http.ResponseWriter, r *http.Request) {
	n := time.Now()
	dt := struct {
		Date string
		Time string
	}{
		n.Format("2006 Jan 02"),
		n.Format("03:04:05 PM"),
	}
	js, err := json.Marshal(dt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// check condition before serving file
func (rh *reqHandler) checkAndServeFile(w http.ResponseWriter, r *http.Request) {
	// handle only GET and POST
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// serve file
	rh.filehandler.ServeHTTP(w, r)
}

// handler for form submit
func (rh *reqHandler) submit(w http.ResponseWriter, r *http.Request) {
	// handle only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Printf("key: %v ==> value: %v\n", k, v)
	}
	http.Redirect(w, r, "https://google.com", 302)
}
