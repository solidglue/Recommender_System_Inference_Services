package server

import "net/http"

type restInferInterface interface {
	restInferServer(w http.ResponseWriter, r *http.Request)
}
