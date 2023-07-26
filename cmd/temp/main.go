package main

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myApp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {

	myHandler := http.HandlerFunc(myApp)

	chain := alice.New(timeoutHandler, nosurf.NewPure).Then()
	http.ListenAndServe(":8000", chain)
}
