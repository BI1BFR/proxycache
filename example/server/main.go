package main

import (
	"net/http"

	"github.com/huangml/proxycache"
)

func main() {
	p := proxycache.New(NewInMemoryDB(), 10, 1, 1)
	http.Handle("/pc/", http.StripPrefix("/pc/", p.HTTPHandlerV1()))
	http.ListenAndServe(":80", nil)
}
