package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// handler funcs for diff req
type reqHandler struct {
	filehandler http.Handler
	userRepo    userRepo
	session     *session
	failures    *attempts
}

func (rh *reqHandler) init() {
	if rh.userRepo != nil {
		rh.userRepo.init()
	}
	rh.session = &session{}
	rh.session.init()
	rh.failures = &attempts{}
	rh.failures.init()
}

func (rh *reqHandler) finalize() {
	if rh.userRepo != nil {
		rh.userRepo.close()
	}
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
	// check if existing session cookie
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
	rdURL := config().RedirectURL
	// check if it was Login OR Cancel action
	if loginAction(r.Form) {
		// check if username & password valid
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if rh.userRepo.validateUser(username, password) {
			// generate login session cookie
			// redirect with auth code to redirect url
			ac := rh.session.authenticate(username)
			rdURL = fmt.Sprintf("%s?code=%s", rdURL, ac)
			//log.Printf("Success; authcode = %s\n", ac)
		} else {
			// if N reattempts failed redirect
			// back to redirect url with afilure code
			// else redirect back to login for N attempts
			a := rh.failures.allowed(username)
			if a {
				// if allowed reattempt
				// redirect to login home
				rdURL = r.Header.Get("Origin")
			}
		}
	}
	// handle Cancel
	// redirect back to 'redirect url'
	http.Redirect(w, r, rdURL, http.StatusTemporaryRedirect)
}

// handler for getting token
func (rh *reqHandler) token(w http.ResponseWriter, r *http.Request) {
	// handle only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// placeholder token requst body
	tr := struct {
		Code string
	}{
		"",
	}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Req. for Token = %v\n", tr.Code)
	// decode tr.code (authcode)
	// extract username and timestamp
	// check in session
	// if valid login session gen id-token
	// set id-token
	// send reponse back with id-token
	// rsp := struct {
	// 	IdToken string
	// }{
	// 	token,
	// }
	// js, err := json.Marshal(rsp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
}

func loginAction(form url.Values) bool {
	isLoginAction := false
	if form.Get("login") != "" {
		isLoginAction = true
	}
	return isLoginAction
}
