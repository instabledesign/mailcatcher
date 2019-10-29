package handler

import (
	"net/http"
)

func Home() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, "nogit/interface/index.html")
		return
	}
}
