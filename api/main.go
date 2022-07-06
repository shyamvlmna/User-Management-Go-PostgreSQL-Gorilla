package main

import (
	"crudapp/router"
	"net/http"
)

func main() {
	r := router.Router()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
