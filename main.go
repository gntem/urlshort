package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

const (
	// HOST address to listen
	HOST = "127.0.0.1"
	// PORT application port
	PORT = "8000"
)

var addr = strings.Join([]string{HOST, PORT}, ":")

// MemStorage hash
var MemStorage = make(map[string]string)

func shortURL(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	urlVar := req.FormValue("url")

	// is valid url
	_, err := url.ParseRequestURI(urlVar)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid url")
	}

	var uniqID = xid.New().String()
	MemStorage[uniqID] = urlVar

	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, "%s/l/%s", addr, uniqID)
}

func redirect(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var httpVariables = mux.Vars(req)

	originalURL, ok := MemStorage[httpVariables["key"]]

	if ok {
		res.WriteHeader(http.StatusTemporaryRedirect)
		http.Redirect(res, req, originalURL, http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "Not found")
	}
}

func listAll(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, MemStorage)
}

func flush(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	totalKeys := len(MemStorage)
	if totalKeys > 1 {
		MemStorage = make(map[string]string)
		fmt.Fprintf(res, "Deleted %d keys", totalKeys)
	} else {
		fmt.Fprintf(res, "Nothing to delete")
	}
}

func main() {

	Router := mux.NewRouter()

	Router.HandleFunc("/short", shortURL).Methods("POST")
	Router.HandleFunc("/l/{key}", redirect).Methods("GET")
	Router.HandleFunc("/listAll", listAll).Methods("GET")
	Router.HandleFunc("/flush", flush).Methods("GET")

	fmt.Println("Serving at :", addr)

	log.Fatal(http.ListenAndServe(addr, Router))

}
