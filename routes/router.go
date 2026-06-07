package routes

import (
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	RegisterAuthRoutes(mux)
	RegisterBookRoutes(mux)
	RegisterTxRoutes(mux)

	return mux
}
