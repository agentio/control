package main

import (
	"flag"
	"fmt"
	"net/http"
)

var port = flag.Uint("p", 8080, "the port to use for serving HTTP requests")

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func main() {
	flag.Parse()
	
	path := fmt.Sprintf(":%v", *port)
	fmt.Printf("%v\n", path)
	var h Hello
	http.ListenAndServe(path, h)
}
