package main

import (
	"crudapp/pkg/router"
	"net/http"
)

func main() {
	r := router.Router()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
