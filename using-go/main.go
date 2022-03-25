package main

import (
	"fmt"
	"net/http"
	"time"
)



type MyFirstHandler string

func (g MyFirstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC3339)
	w.Write([]byte("The time is: " + tm))
}

func main() {
	var firstHandler MyFirstHandler
	// declare our serveMux in main
	router :=http.NewServeMux()
	// register handler we defined - it now responds to any request to path use-handler
	router.Handle("/use-handler", firstHandler)

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

