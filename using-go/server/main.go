package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type MyFirstHandler struct{}

func (g MyFirstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	tm := time.Now().Format(time.RFC3339)
	w.Write([]byte(fmt.Sprintf("hello %s, the time is: %s\n", name, tm)))
}

func main() {
	var firstHandler MyFirstHandler
	// declare our serveMux in main
	router := mux.NewRouter()
	// register handler we defined - it now responds to any request to path use-handler
	router.Handle("/use-handler/{name:[a-zA-Z]+}", firstHandler).Methods("GET")
	router.Handle("/heath-check",  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("OK"))
		return
	}))

	// ListenAndServe listens on the TCP network address addr and then calls
	// Serve with handler to handle requests on incoming connections.
	// Accepted connections are configured to enable TCP keep-alives.
	//
	// The handler is typically nil, in which case the DefaultServeMux is used.
	//
	// ListenAndServe always returns a non-nil error.
	//ListenAndServer func signature -- ListenAndServe(addr string, handler Handler) error
	err := http.ListenAndServe(":8833", router)
	if err != nil {
		fmt.Println("error while attempting to listen for incoming connections", err)
	}
}
